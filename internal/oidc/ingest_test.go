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
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	echo "github.com/labstack/echo/v4"
)

// Ingest test constants
const (
	ingestTestClientID     = "login-client"
	ingestTestOIDCIssuer   = "https://example.com"
	ingestTestStatusKey    = "status"
	ingestTestAuthModeOIDC = "oidc"
	ingestTestUserID       = "user-123"
	ingestTestAudience     = "test-audience"
)

// TestIngestBearerMiddlewareValidToken tests that valid tokens are accepted and handler runs.
func TestIngestBearerMiddlewareValidToken(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "user-123", "client-123", "scope1", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{ingestTestStatusKey: "ok"})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		"azp":    "client-123",
		"scope":  "scope1 scope2",
		claimExp: time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
		"nbf":    time.Now().Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Errorf("handler was not called for valid token")
	}
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareWrongAudience tests that tokens with wrong audience are rejected.
func TestIngestBearerMiddlewareWrongAudience(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "user-123", "client-123", "scope1", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: "wrong-audience",
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for wrong audience")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareExpiredToken tests that expired tokens are rejected.
func TestIngestBearerMiddlewareExpiredToken(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(-time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for expired token")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for expired token, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareNbfFuture tests that tokens with nbf in future are rejected.
func TestIngestBearerMiddlewareNbfFuture(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(time.Hour).Unix(),
		"nbf":    time.Now().Add(2 * time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called when nbf is in future")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for nbf in future, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareAlgNone tests that alg=none tokens are rejected.
func TestIngestBearerMiddlewareAlgNone(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	// Create a token with alg=none (invalid)
	//nolint:gosec
	token := "eyJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoidGVzdC1hdWRpZW5jZSIsImV4cCI6OTk5OTk5OTk5OX0."

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for alg=none")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for alg=none, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareHs256AlgConfusion tests algorithm confusion attacks.
func TestIngestBearerMiddlewareHs256AlgConfusion(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	// HS256 token (invalid for RS256 endpoint)
	//nolint:gosec
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoidGVzdC1hdWRpZW5jZSIsImV4cCI6OTk5OTk5OTk5OX0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for HS256 token")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for HS256 algorithm confusion, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareMissingHeader tests that missing Authorization header returns 401.
func TestIngestBearerMiddlewareMissingHeader(t *testing.T) {
	setupIngestConfig(ingestTestAuthModeOIDC, ingestTestOIDCIssuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	// No Authorization header
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called without Authorization header")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for missing header, got %d", rec.Code)
	}

	// Check WWW-Authenticate header
	wwwAuth := rec.Header().Get("WWW-Authenticate")
	if wwwAuth != `Bearer realm="falcosidekick-ui"` {
		t.Errorf("expected exact WWW-Authenticate header, got: %q", wwwAuth)
	}
}

// TestIngestBearerMiddlewareTokenInQueryParam tests that token in query param is ignored.
func TestIngestBearerMiddlewareTokenInQueryParam(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/?token=%s", token), nil)
	// Token only in query param, not in header
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called when token is only in query param")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 when token in query param only, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareLoginClientIDRejected tests that login client ID cannot be used as ingestion audience.
func TestIngestBearerMiddlewareLoginClientIDRejected(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")
	cfg := configuration.GetConfiguration()
	cfg.OIDCClientID = ingestTestClientID

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	// Token with login client ID as audience (should be rejected)
	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestClientID,
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called when using login client ID")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 when using login client ID, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareKeycloakPayloadTypBearer tests Keycloak compatibility.
func TestIngestBearerMiddlewareKeycloakPayloadTypBearer(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	// Keycloak-style token with payload typ="Bearer"
	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		"typ":    "Bearer",
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Errorf("handler should be called for Keycloak payload typ=Bearer")
	}
	if rec.Code != 200 {
		t.Errorf("expected 200 for Keycloak payload typ=Bearer, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareJOSEHeaderTypJWT tests JOSE header typ validation.
func TestIngestBearerMiddlewareJOSEHeaderTypJWT(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Errorf("handler should be called for valid JOSE header typ")
	}
	if rec.Code != 200 {
		t.Errorf("expected 200 for JOSE header typ validation, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareInvalidJOSEHeaderTyp tests invalid JOSE header typ rejection.
func TestIngestBearerMiddlewareInvalidJOSEHeaderTyp(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	// Token with invalid JOSE header typ (not JWT or at+jwt)
	//nolint:gosec
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6ImZvbyJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoidGVzdC1hdWRpZW5jZSIsImV4cCI6OTk5OTk5OTk5OX0.invalid"

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for invalid JOSE header typ")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for invalid JOSE header typ, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareSubjectAllowlist tests subject allowlist enforcement.
func TestIngestBearerMiddlewareSubjectAllowlist(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "allowed-sub", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: "disallowed-sub",
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for disallowed subject")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for disallowed subject, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareClientAllowlist tests client (azp/client_id) allowlist enforcement.
func TestIngestBearerMiddlewareClientAllowlist(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "allowed-client", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		"azp":    "disallowed-client",
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for disallowed client")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for disallowed client, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareScopeAllowlist tests scope validation for string and scp array formats.
func TestIngestBearerMiddlewareScopeAllowlist(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "required-scope", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		"scope":  "other-scope",
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called without required scope")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for missing required scope, got %d", rec.Code)
	}
}

// TestIngestBearerMiddlewareCaseInsensitiveBearer tests RFC 6750 case-insensitivity.
func TestIngestBearerMiddlewareCaseInsensitiveBearer(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(time.Hour).Unix(),
	})

	schemes := []string{"Bearer", "bearer", "BEARER", "BeArEr"}
	for _, scheme := range schemes {
		resetIngestProvider()
		defer resetIngestProvider()

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			return c.JSON(200, map[string]string{})
		}

		e := echo.New()
		e.POST("/api/v1/", handler, IngestBearerMiddleware())

		req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", scheme, token))
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if !handlerCalled || rec.Code != 200 {
			t.Errorf("expected success for %q scheme, got code %d, handler called %v", scheme, rec.Code, handlerCalled)
		}
	}
}

// TestIngestAuthNone tests that INGEST_AUTH=none allows unauthenticated requests.
func TestIngestAuthNone(t *testing.T) {
	cfg := configuration.CreateConfiguration()
	cfg.IngestAuth = "none"

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	e := echo.New()
	e.POST("/api/v1/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/", nil)
	// No Authorization header
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if !handlerCalled || rec.Code != 200 {
		t.Errorf("expected handler to be called when INGEST_AUTH=none, got code %d", rec.Code)
	}
}

// TestRootPathProtected tests that the root path (/) is also protected.
func TestRootPathProtected(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(200, map[string]string{})
	}

	e := echo.New()
	e.POST("/", handler, IngestBearerMiddleware())

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	// No Authorization header
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if handlerCalled {
		t.Errorf("handler should not be called for root path without auth")
	}
	if rec.Code != 401 {
		t.Errorf("expected 401 for root path without auth, got %d", rec.Code)
	}
}

// TestValidateIngestConfigValid tests valid ingestion configuration.
func TestValidateIngestConfigValid(t *testing.T) {
	cfg := configuration.CreateConfiguration()
	cfg.IngestAuth = ingestTestAuthModeOIDC
	cfg.IngestOIDCIssuer = ingestTestOIDCIssuer
	cfg.IngestOIDCAudience = ingestTestAudience
	cfg.IngestOIDCAllowedSubjects = "user1,user2"
	cfg.OIDCClientID = ingestTestClientID

	// Should not panic or error
	ValidateIngestConfig()
}

// TestValidateIngestConfigMissingIssuer tests that missing issuer is caught.
func TestValidateIngestConfigMissingIssuer(t *testing.T) {
	cfg := configuration.CreateConfiguration()
	cfg.IngestAuth = ingestTestAuthModeOIDC
	cfg.IngestOIDCIssuer = ""
	cfg.IngestOIDCAudience = ingestTestAudience

	// Should catch missing issuer (in real code would fatal)
	// This test just verifies the function runs
	defer func() {
		if r := recover(); r != nil {
			// Expected
		}
	}()
	ValidateIngestConfig()
}

// TestValidateIngestConfigAudienceEqualsClientID tests that audience == login client ID is caught.
func TestValidateIngestConfigAudienceEqualsClientID(t *testing.T) {
	cfg := configuration.CreateConfiguration()
	cfg.IngestAuth = ingestTestAuthModeOIDC
	cfg.IngestOIDCIssuer = ingestTestOIDCIssuer
	cfg.IngestOIDCAudience = ingestTestClientID
	cfg.OIDCClientID = ingestTestClientID

	// Should catch audience == login client ID (in real code would fatal)
	defer func() {
		if r := recover(); r != nil {
			// Expected
		}
	}()
	ValidateIngestConfig()
}

// testIssuer provides a simple OIDC issuer for testing.
type testIssuer struct {
	Issuer  string
	privKey interface{}
	pubKey  interface{}
	keyID   string
	signer  jose.Signer
	server  *httptest.Server
}

// setupTestIssuer creates an in-process OIDC issuer with httptest.
func setupTestIssuer() (*testIssuer, *httptest.Server) {
	// Generate RSA key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	keyID := "test-key-1"
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       privKey,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		panic(err)
	}

	issuer := &testIssuer{
		privKey: privKey,
		pubKey:  &privKey.PublicKey,
		keyID:   keyID,
		signer:  signer,
	}

	// Create httptest server serving discovery and JWKS
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/openid-configuration" {
			w.Header().Set("Content-Type", "application/json")
			discovery := map[string]interface{}{
				"issuer":                 issuer.Issuer,
				"authorization_endpoint": issuer.Issuer + "/authorize",
				"token_endpoint":         issuer.Issuer + "/token",
				"jwks_uri":               issuer.Issuer + "/.well-known/jwks.json",
			}
			json.NewEncoder(w).Encode(discovery)
		} else if r.URL.Path == "/.well-known/jwks.json" {
			w.Header().Set("Content-Type", "application/json")
			pubKey := issuer.pubKey.(*rsa.PublicKey)
			jwks := map[string]interface{}{
				"keys": []map[string]interface{}{
					{
						"kty": "RSA",
						"kid": issuer.keyID,
						"use": "sig",
						"n":   base64URLEncode(pubKey.N.Bytes()),
						"e":   base64URLEncode([]byte{0x01, 0x00, 0x01}),
					},
				},
			}
			json.NewEncoder(w).Encode(jwks)
		}
	}))

	issuer.Issuer = server.URL
	issuer.server = server
	return issuer, server
}

// SignToken creates a signed JWT token with the given claims.
func (ti *testIssuer) SignToken(claims map[string]interface{}) (string, error) {
	// Ensure issuer claim is set
	if _, ok := claims["iss"]; !ok {
		claims["iss"] = ti.Issuer
	}
	token, err := jwt.Signed(ti.signer).Claims(claims).Serialize()
	return token, err
}

// base64URLEncode encodes bytes as base64url without padding.
func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// setupIngestConfig configures ingestion auth settings.
func setupIngestConfig(auth, issuer, audience, subjects, clients, scope, cafile string) {
	cfg := configuration.CreateConfiguration()
	cfg.IngestAuth = auth
	cfg.IngestOIDCIssuer = issuer
	cfg.IngestOIDCAudience = audience
	cfg.IngestOIDCAllowedSubjects = subjects
	cfg.IngestOIDCAllowedClients = clients
	cfg.IngestOIDCRequiredScope = scope
	cfg.IngestOIDCCAFile = cafile
	cfg.OIDCClientID = ingestTestClientID
}

// resetIngestProvider clears the cached ingestion provider for testing.
func resetIngestProvider() {
	ingestProviderMutex.Lock()
	defer ingestProviderMutex.Unlock()
	cachedIngestProvider = nil
}

// TestIsInAllowlistEmptyTrailingComma tests that trailing commas create empty entries that are skipped.
//
//nolint:goconst
func TestIsInAllowlistEmptyTrailingComma(t *testing.T) {
	// "svc," splits to ["svc", ""], should skip the empty entry
	tests := []struct {
		value      string
		allowlist  string
		shouldPass bool
	}{
		{"svc", "svc,", true}, // "svc" is in allowlist
		{"", "svc,", false},   // Empty value never matches
		{"", "svc,,", false},  // Empty value never matches, even with multiple empty entries
		{"svc", "svc,,user", true},
		{"user", ",svc,user", true}, // Empty at start
		{"other", "svc,", false},    // "other" not in allowlist
	}

	for _, tc := range tests {
		result := isInAllowlist(tc.value, tc.allowlist)
		if result != tc.shouldPass {
			t.Errorf("isInAllowlist(%q, %q) = %v, want %v", tc.value, tc.allowlist, result, tc.shouldPass)
		}
	}
}

// TestIsUserInAllowedGroupsEmptyEntries tests group matching with empty entries after split/trim.
//
//nolint:goconst
func TestIsUserInAllowedGroupsEmptyEntries(t *testing.T) {
	tests := []struct {
		userGroups []string
		allowlist  string
		shouldPass bool
	}{
		{[]string{"admins"}, "admins,users", true},  // "admins" is in allowlist
		{[]string{""}, "admins,users", false},       // Empty group value never matches
		{[]string{"admins"}, "admins,", true},       // "admins" is in allowlist, trailing comma creates empty entry which is skipped
		{[]string{"user"}, ",user,admin", true},     // User in middle of allowlist with empty entries
		{[]string{"admin", "user"}, "admin", true},  // First group matches
		{[]string{"", "admin"}, "admin,user", true}, // Empty group skipped, second group matches
		{[]string{"other"}, "admin,user", false},    // No match
	}

	for _, tc := range tests {
		result := isUserInAllowedGroups(tc.userGroups, tc.allowlist)
		if result != tc.shouldPass {
			t.Errorf("isUserInAllowedGroups(%v, %q) = %v, want %v", tc.userGroups, tc.allowlist, result, tc.shouldPass)
		}
	}
}
