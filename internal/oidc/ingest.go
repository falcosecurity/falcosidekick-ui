// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

// IngestBearerMiddleware returns middleware that validates bearer tokens for ingestion endpoints.
// RFC 6750: Authorization header with Bearer scheme.
// Tokens are validated against the ingestion OIDC issuer and audience.
// This middleware should only be attached to ingestion routes (route-level middleware, not global).
func IngestBearerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			config := configuration.GetConfiguration()

			// Only validate if ingest auth is enabled
			if config.IngestAuth != authModeOIDC {
				return next(c)
			}

			// Extract bearer token from Authorization header (RFC 6750)
			// RFC 6750: auth-scheme is case-insensitive
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// Missing header: 401 with basic WWW-Authenticate
				c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui"`)
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
			}

			// Parse "Bearer <token>" format (case-insensitive scheme)
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				// Invalid format
				c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui", error="invalid_request"`)
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
			}

			token := strings.TrimSpace(parts[1])
			if token == "" {
				// Empty token
				c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui", error="invalid_request"`)
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
			}

			// Validate token before reading body
			claims, err := validateIngestToken(c.Request().Context(), token)
			if err != nil {
				utils.WriteLog("warning", fmt.Sprintf("ingestion token validation failed: %v", err))
				c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui", error="invalid_token"`)
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
			}

			// Check subject allowlist
			if config.IngestOIDCAllowedSubjects != "" {
				sub, ok := claims[claimSub].(string)
				if !ok || !isInAllowlist(sub, config.IngestOIDCAllowedSubjects) {
					utils.WriteLog("warning", "ingestion token subject not in allowed list")
					c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui", error="invalid_token"`)
					return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
				}
			}

			// Check client allowlist (azp or client_id)
			if config.IngestOIDCAllowedClients != "" {
				azp, _ := claims[claimAzp].(string)
				clientID, _ := claims[claimClientID].(string)
				client := azp
				if client == "" {
					client = clientID
				}
				if client == "" || !isInAllowlist(client, config.IngestOIDCAllowedClients) {
					utils.WriteLog("warning", "ingestion token client not in allowed list")
					c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui", error="invalid_token"`)
					return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
				}
			}

			// Check required scope
			if config.IngestOIDCRequiredScope != "" {
				if !hasRequiredScope(claims, config.IngestOIDCRequiredScope) {
					utils.WriteLog("warning", "ingestion token missing required scope")
					c.Response().Header().Set("WWW-Authenticate", `Bearer realm="falcosidekick-ui", error="invalid_token"`)
					return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthorized})
				}
			}

			// Token is valid, proceed to handler
			return next(c)
		}
	}
}

// validateIngestToken validates a bearer token against the ingestion OIDC provider.
func validateIngestToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	config := configuration.GetConfiguration()

	// Get ingest provider (separate from login provider)
	provider, err := getIngestProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingestion OIDC provider: %w", err)
	}

	// Create verifier for the ingestion audience
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.IngestOIDCAudience,
		SupportedSigningAlgs: []string{
			"RS256", "RS384", "RS512",
			"ES256", "ES384", "ES512",
			"PS256", "PS384", "PS512",
		},
	})

	// Verify token
	idToken, err := verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	// Extract claims
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	// Validation: JOSE header typ (RFC 7519, Section 5.1)
	// Check the JWT header typ field (case-insensitive), not the payload typ.
	// Valid values: "JWT" or "at+jwt" (RFC 9068). If typ is absent, that's OK.
	if headerTyp, err := getJOSEHeaderTyp(tokenString); err == nil && headerTyp != "" {
		// typ is present in JOSE header - must be "JWT" or "at+jwt" (case-insensitive)
		if !strings.EqualFold(headerTyp, "JWT") && !strings.EqualFold(headerTyp, "at+jwt") {
			return nil, fmt.Errorf("invalid JOSE header typ: %q", headerTyp)
		}
	}

	// Defense in depth: reject Keycloak ID tokens (payload typ == "ID")
	// go-oidc validates aud, so this is extra protection against misconfiguration.
	if payloadTyp, ok := claims["typ"].(string); ok && strings.EqualFold(payloadTyp, "ID") {
		return nil, fmt.Errorf("ID token not allowed for ingestion (use access token)")
	}

	// nbf (not-before) validation:
	// go-oidc's verifier already enforces exp and nbf by default (with clock skew).
	// We trust go-oidc's validation and do not duplicate it here.

	return claims, nil
}

// getJOSEHeaderTyp extracts the JOSE header "typ" field from a JWT.
// JWT format: base64url(header).base64url(payload).base64url(signature)
// Returns empty string if typ is not present or cannot be extracted.
func getJOSEHeaderTyp(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid JWT format")
	}

	// Decode the header (first part)
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT header: %w", err)
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return "", fmt.Errorf("failed to parse JWT header: %w", err)
	}

	// Extract typ field (case-sensitive key name, but value comparison is case-insensitive)
	if typ, ok := header["typ"].(string); ok {
		return typ, nil
	}

	// typ field not present in header - this is allowed per RFC 7519
	return "", nil
}

// isInAllowlist checks if a value is in a comma-separated allowlist.
// Rejects empty entries and never matches an empty value.
func isInAllowlist(value, allowlistStr string) bool {
	// Never match empty value
	if value == "" {
		return false
	}
	allowlist := strings.Split(allowlistStr, ",")
	for _, item := range allowlist {
		trimmed := strings.TrimSpace(item)
		// Skip empty entries after trim
		if trimmed == "" {
			continue
		}
		if trimmed == value {
			return true
		}
	}
	return false
}

// hasRequiredScope checks if a required scope is present in the token
func hasRequiredScope(claims map[string]interface{}, requiredScope string) bool {
	// Check scope claim (space-separated string)
	if scope, ok := claims["scope"].(string); ok {
		scopes := strings.Fields(scope)
		for _, s := range scopes {
			if s == requiredScope {
				return true
			}
		}
	}

	// Check scp claim (array)
	if scp, ok := claims["scp"].([]interface{}); ok {
		for _, s := range scp {
			if scpStr, ok := s.(string); ok && scpStr == requiredScope {
				return true
			}
		}
	}

	return false
}
