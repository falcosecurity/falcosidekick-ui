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
	"errors"
	"net/http"
	"time"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	echo "github.com/labstack/echo/v4"
)

// HTTP method constants
const (
	methodGET  = "GET"
	methodPOST = "POST"
)

// API endpoint constants
const (
	pathHealthz      = "/api/v1/healthz"
	pathAuthMe       = "/api/v1/auth/me"
	pathLoginOIDC    = "/api/v1/auth/oidc/login"
	pathCallbackOIDC = "/api/v1/auth/oidc/callback"
	pathIngestRoot   = "/api/v1/"
	pathIngestAdd    = "/api/v1/events/add"
	pathLogout       = "/api/v1/auth/logout"
	pathAuth         = "/api/v1/auth"
	pathAuthenticate = "/api/v1/authenticate"
)

// SessionMiddleware returns an echo middleware that validates OIDC sessions.
// Returns a middleware factory function with the proper echo signature.
func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			config := configuration.GetConfiguration()
			path := c.Path()
			method := c.Request().Method

			// Skip if not in OIDC mode
			if config.AuthMode != authModeOIDC {
				return next(c)
			}

			// Public endpoints (no session required)
			if isPublicEndpoint(path, method) {
				return next(c)
			}

			// Get session cookie
			sessionCookie, err := c.Cookie("__Host-fsui_session")
			if err != nil && config.OIDCInsecureAllowHTTP {
				sessionCookie, err = c.Cookie("fsui_session")
			}

			if err != nil {
				c.Response().Header().Set("Content-Type", "application/json")
				c.Response().Header().Set("Cache-Control", "no-store")
				c.Response().Header().Set("Pragma", "no-cache")
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthenticated})
			}

			// Get session from Redis
			conn, err := GetRedisConn()
			if err != nil {
				c.Response().Header().Set("Content-Type", "application/json")
				c.Response().Header().Set("Cache-Control", "no-store")
				c.Response().Header().Set("Pragma", "no-cache")
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthenticated})
			}
			defer conn.Close()

			session, err := GetSession(conn, sessionCookie.Value)
			if err != nil {
				c.Response().Header().Set("Content-Type", "application/json")
				c.Response().Header().Set("Cache-Control", "no-store")
				c.Response().Header().Set("Pragma", "no-cache")
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthenticated})
			}

			// Check absolute TTL
			createdAt := session.CreatedAt
			ttl := time.Duration(config.SessionTTL) * time.Second

			if time.Since(createdAt) >= ttl {
				DeleteSession(conn, sessionCookie.Value)
				c.Response().Header().Set("Content-Type", "application/json")
				c.Response().Header().Set("Cache-Control", "no-store")
				c.Response().Header().Set("Pragma", "no-cache")
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthenticated})
			}

			// Refresh idle timeout
			if err := RefreshSessionTTL(conn, sessionCookie.Value); err != nil {
				// Only delete on ErrSessionExpired; on transient Redis errors, keep session
				if errors.Is(err, ErrSessionExpired) {
					DeleteSession(conn, sessionCookie.Value)
				}
				// Return 401 on any error (session expired or transient Redis error)
				c.Response().Header().Set("Content-Type", "application/json")
				c.Response().Header().Set("Cache-Control", "no-store")
				c.Response().Header().Set("Pragma", "no-cache")
				return c.JSON(http.StatusUnauthorized, map[string]string{ErrorJSON: errorUnauthenticated})
			}

			// Store username in context
			c.Set("username", session.Username)

			// Session is valid; proceed to next handler
			return next(c)
		}
	}
}

// isPublicEndpoint checks if the endpoint is public.
func isPublicEndpoint(path string, method string) bool {
	// Public GET endpoints
	if method == methodGET {
		if path == pathHealthz {
			return true
		}
		if path == pathAuthMe {
			return true
		}
		if path == pathLoginOIDC {
			return true
		}
		if path == pathCallbackOIDC {
			return true
		}
	}

	// Public POST endpoints (event ingestion)
	if method == methodPOST {
		if path == pathIngestRoot {
			return true
		}
		if path == pathIngestAdd {
			return true
		}
		if path == pathLogout {
			return true
		}
		// In OIDC mode, basic auth endpoints return 404 (not 401) since they don't exist
		if path == pathAuth {
			return true
		}
		if path == pathAuthenticate {
			return true
		}
	}

	return false
}
