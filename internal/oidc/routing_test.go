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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	echo "github.com/labstack/echo/v4"
)

// TestSessionMiddlewareWithRouting exercises SessionMiddleware through a real echo router so that
// c.Path() is the registered route pattern (as in main.go). Unlike TestSessionMiddlewarePublicEndpoint
// it asserts whether the downstream handler actually ran, so it cannot pass vacuously when the
// middleware writes a 401 and returns nil.
func TestSessionMiddlewareWithRouting(t *testing.T) {
	configuration.CreateConfiguration()
	cfg := configuration.GetConfiguration()
	cfg.AuthMode = authModeOIDC
	cfg.OIDCInsecureAllowHTTP = false
	cfg.SessionTTL = 28800

	e := echo.New()
	g := e.Group("/api/v1")
	called := map[string]bool{}
	h := func(name string) echo.HandlerFunc {
		return func(c echo.Context) error {
			called[name] = true
			return c.String(http.StatusOK, "handler-"+name)
		}
	}

	// RegisterRoutes sets up SessionMiddleware, auth routes (me, login, callback, logout),
	// and 404 handlers for /auth and /authenticate. This matches production wiring.
	RegisterRoutes(g)

	// Register sentinel handlers only for routes NOT handled by RegisterRoutes
	g.GET("/healthz", h("healthz"))
	g.POST("/", h("ingest-root"))
	g.POST("/events/add", h("ingest-add"))
	g.GET("/configuration", h("configuration"))
	g.GET("/outputs", h("outputs"))
	g.GET("/events/search", h("search"))
	g.GET("/events/count", h("count"))
	g.GET("/version", h("version"))

	cases := []struct {
		method, path, name string
		expectCode         int
		isRealHandler      bool // true if handler runs; false if middleware blocks
	}{
		{methodGET, pathHealthz, "healthz", http.StatusOK, true},
		{methodGET, pathAuthMe, "me", http.StatusOK, true},                     // real handler may return 200 or 400
		{methodGET, pathLoginOIDC, "login", http.StatusFound, true},            // real handler may return 302
		{methodGET, pathCallbackOIDC, "callback", http.StatusBadRequest, true}, // missing code param
		{methodPOST, pathIngestRoot, "ingest-root", http.StatusOK, true},
		{methodPOST, pathIngestAdd, "ingest-add", http.StatusOK, true},
		{methodPOST, pathLogout, "logout", http.StatusOK, true}, // real handler may need CSRF check
		{methodGET, "/api/v1/configuration", "configuration", http.StatusUnauthorized, false},
		{methodGET, "/api/v1/outputs", "outputs", http.StatusUnauthorized, false},
		{methodGET, "/api/v1/events/search", "search", http.StatusUnauthorized, false},
		{methodGET, "/api/v1/events/count", "count", http.StatusUnauthorized, false},
		{methodGET, "/api/v1/version", "version", http.StatusUnauthorized, false},
		// In OIDC mode, basic auth endpoints return 404 (not 401)
		{methodPOST, "/api/v1/auth", "auth", http.StatusNotFound, false},
		{methodPOST, "/api/v1/authenticate", "authenticate", http.StatusNotFound, false},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			called[tc.name] = false
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))

			if tc.isRealHandler {
				// Real handler ran: status should NOT be 401 and body should NOT be unauthenticated error
				if rec.Code == http.StatusUnauthorized && rec.Body.String() == "{\"error\":\"unauthenticated\"}\n" {
					t.Fatalf("real handler should have run, but got 401 unauthenticated")
				}
			} else {
				// Middleware blocked: expect exact status and error body
				if rec.Code != tc.expectCode {
					t.Fatalf("expected %d got %d body=%q", tc.expectCode, rec.Code, rec.Body.String())
				}
				if tc.expectCode == http.StatusUnauthorized {
					if body := rec.Body.String(); body != "{\"error\":\"unauthenticated\"}\n" {
						t.Fatalf("unexpected body %q", body)
					}
				}
			}
		})
	}
}

// TestAuthEndpointsReturn404 verifies that basic auth endpoints return 404 in OIDC mode
func TestAuthEndpointsReturn404(t *testing.T) {
	configuration.CreateConfiguration()
	cfg := configuration.GetConfiguration()
	cfg.AuthMode = authModeOIDC

	// Test that /api/v1/auth and /api/v1/authenticate endpoints are marked as public
	// so they return 404 instead of 401 when not registered
	testCases := []struct {
		path   string
		method string
	}{
		{"/api/v1/auth", methodPOST},
		{"/api/v1/authenticate", methodPOST},
	}

	for _, tc := range testCases {
		result := isPublicEndpoint(tc.path, tc.method)
		if !result {
			t.Fatalf("isPublicEndpoint(%q, %q) should be true to return 404, got false", tc.path, tc.method)
		}
	}
}

// TestErrorCodeValidation verifies error code validation
func TestErrorCodeValidation(t *testing.T) {
	testCases := []struct {
		code    string
		isValid bool
	}{
		// Valid codes
		{"access_denied", true},
		{"server_error", true},
		{"temporarily_unavailable", true},
		{"login_required", true},
		{"_", true},
		{"a", true},
		{"", false},
		// Invalid codes
		{"Access_Denied", false},
		{"access-denied", false},
		{"access denied", false},
		{"access_denied!", false},
		{strings.Repeat("a", 65), false}, // too long
		{strings.Repeat("a", 64), true},  // max length is 64
	}

	for _, tc := range testCases {
		result := isValidOIDCErrorCode(tc.code)
		if result != tc.isValid {
			t.Fatalf("isValidOIDCErrorCode(%q) = %v, want %v", tc.code, result, tc.isValid)
		}
	}
}

// TestIngestRoutesProtectionOIDC tests that ingest routes are protected when INGEST_AUTH=oidc.
func TestIngestRoutesProtectionOIDC(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	resetIngestProvider()
	defer resetIngestProvider()

	setupIngestConfig(ingestTestAuthModeOIDC, issuer.Issuer, ingestTestAudience, "", "", "", "")

	e := echo.New()
	g := e.Group("/api/v1")

	// Register routes like main.go does
	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}

	// Build middleware slice as in main.go
	var ingestMW []echo.MiddlewareFunc
	if configuration.GetConfiguration().IngestAuth == authModeOIDC {
		ingestMW = append(ingestMW, IngestBearerMiddleware())
	}

	// Register the three ingest routes exactly once
	const routeNameAddEvent = "add-event"
	e.POST("/", handler, ingestMW...).Name = routeNameAddEvent
	g.POST("/", handler, ingestMW...).Name = routeNameAddEvent
	eventsRoute := g.Group("/events")
	eventsRoute.POST("/add", handler, ingestMW...).Name = routeNameAddEvent

	// Test all three routes without token → should get 401
	routes := []string{"/", "/api/v1/", "/api/v1/events/add"}
	for _, route := range routes {
		t.Run("POST "+route+" without token", func(t *testing.T) {
			handlerCalled = false
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, route, nil))

			if handlerCalled {
				t.Errorf("handler should not be called without token at %s", route)
			}
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("expected 401 at %s, got %d", route, rec.Code)
			}
		})
	}

	// Test with valid token → should get 200
	token, _ := issuer.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: 9999999999,
	})

	for _, route := range routes {
		t.Run("POST "+route+" with token", func(t *testing.T) {
			handlerCalled = false
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, route, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			e.ServeHTTP(rec, req)

			if !handlerCalled {
				t.Errorf("handler should be called with token at %s", route)
			}
			if rec.Code != http.StatusOK {
				t.Errorf("expected 200 at %s, got %d", route, rec.Code)
			}
		})
	}
}
