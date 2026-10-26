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

// Regression tests for logout without a session cookie, JWKS bearer delivery and provider
// initialisation races. Run with -race to exercise the race test.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	echo "github.com/labstack/echo/v4"
)

// Logout must not dereference a nil *http.Cookie when the request carries no session cookie
// (handlers.go: expireCookie uses sessionCookie.Name). Spec: logout must work even if the session expired.
func TestLogoutWithoutSessionCookieDoesNotPanic(t *testing.T) {
	configuration.CreateConfiguration()
	cfg := configuration.GetConfiguration()
	cfg.AuthMode = "oidc"
	cfg.OIDCRedirectURL = "https://ui.example.test/api/v1/auth/oidc/callback"
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://ui.example.test")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Logout panicked without a session cookie: %v", r)
		}
	}()
	if err := Logout(c); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("want 200, got err=%v code=%d", err, rec.Code)
	}
}

// INGEST_OIDC_JWKS_BEARER_FILE must be attached to discovery/JWKS requests; ProxyConnectHeader is only
// sent to an HTTP proxy on CONNECT and would never reach the issuer.
func TestJWKSBearerFileIsSentToIssuer(t *testing.T) {
	var mu sync.Mutex
	var seen []string
	issuer, server := setupTestIssuer()
	defer server.Close()
	// wrap: record Authorization headers on the issuer endpoints
	orig := server.Config.Handler
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.Header.Get("Authorization"))
		mu.Unlock()
		orig.ServeHTTP(w, r)
	})
	dir := t.TempDir()
	f := filepath.Join(dir, "bearer")
	if err := os.WriteFile(f, []byte("sa-token-123\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetIngestProvider()
	defer resetIngestProvider()
	setupIngestConfig("oidc", issuer.Issuer, "aud", "", "", "scope-x", "")
	cfg := configuration.GetConfiguration()
	cfg.IngestOIDCJWKSBearerFile = f
	cfg.IngestOIDCInsecureAllowHTTP = true // Test issuer uses HTTP
	if _, err := getIngestProvider(t.Context()); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, a := range seen {
		if a == "Bearer sa-token-123" {
			return
		}
	}
	t.Fatalf("bearer token from JWKS_BEARER_FILE was never sent to the issuer; Authorization headers seen: %q", seen)
}

// The lazily-initialised provider/verifier/oauth2 globals must be guarded; run with -race.
func TestProviderInitRace(t *testing.T) {
	issuer, server := setupTestIssuer()
	defer server.Close()
	setupTestConfig("oidc", issuer.Issuer, "client")
	configuration.GetConfiguration().OIDCClientSecret = "s"
	providerMutex.Lock()
	cachedProvider, oauth2Config, idTokenVerifier = nil, nil, nil
	providerMutex.Unlock()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = GetOAuth2Config(t.Context())
			_, _ = GetIDTokenVerifier(t.Context())
		}()
	}
	wg.Wait()
}
