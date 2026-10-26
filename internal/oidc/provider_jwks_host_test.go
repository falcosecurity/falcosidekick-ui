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
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	gooidc "github.com/coreos/go-oidc/v3/oidc"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
)

// TestIngestJWKSBearerSplitHost reproduces the Kubernetes layout: the issuer and the
// JWKS live on different hosts and the JWKS host rejects anonymous requests.
func TestIngestJWKSBearerSplitHost(t *testing.T) {
	ti, tiServer := setupTestIssuer()
	defer tiServer.Close()
	pub := ti.pubKey.(*rsa.PublicKey)

	const bearer = "sa-token"
	var jwksCalls, jwksUnauth int32

	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&jwksCalls, 1)
		if r.Header.Get("Authorization") != "Bearer "+bearer {
			atomic.AddInt32(&jwksUnauth, 1)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"keys": []map[string]interface{}{{
				"kty": "RSA", "kid": ti.keyID, "use": "sig", "alg": gooidc.RS256,
				"n": base64URLEncode(pub.N.Bytes()),
				"e": base64URLEncode([]byte{0x01, 0x00, 0x01}),
			}},
		})
	}))
	defer jwksServer.Close()

	var issuerSawBearer int32
	var issuerURL string
	issuerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			atomic.AddInt32(&issuerSawBearer, 1)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"issuer":                 issuerURL,
			"authorization_endpoint": issuerURL + "/authorize",
			"token_endpoint":         issuerURL + "/token",
			"jwks_uri":               jwksServer.URL + "/openid/v1/jwks",
		})
	}))
	defer issuerServer.Close()
	issuerURL = issuerServer.URL
	ti.Issuer = issuerURL

	tokenFile := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tokenFile, []byte(bearer+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetIngestProvider()
	defer resetIngestProvider()
	setupIngestConfig(ingestTestAuthModeOIDC, issuerURL, ingestTestAudience, "", "", "", "")
	cfg := configuration.GetConfiguration()
	cfg.IngestOIDCInsecureAllowHTTP = true
	cfg.IngestOIDCJWKSBearerFile = tokenFile

	token, err := ti.SignToken(map[string]interface{}{
		claimSub: ingestTestUserID,
		claimAud: ingestTestAudience,
		claimExp: time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := validateIngestToken(context.Background(), token); err != nil {
		t.Fatalf("token should validate when JWKS requires the bearer: %v (jwks calls=%d, unauthenticated=%d)",
			err, atomic.LoadInt32(&jwksCalls), atomic.LoadInt32(&jwksUnauth))
	}
	if atomic.LoadInt32(&jwksUnauth) != 0 {
		t.Errorf("JWKS host received %d unauthenticated requests", jwksUnauth)
	}
}

// TestBearerRoundTripperHostScoping checks the token only goes to trusted hosts over HTTPS.
func TestBearerRoundTripperHostScoping(t *testing.T) {
	tokenFile := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tokenFile, []byte("tok"), 0o600); err != nil {
		t.Fatal(err)
	}
	var got string
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		got = r.Header.Get("Authorization")
		return &http.Response{StatusCode: 200, Body: http.NoBody, Request: r}, nil
	})
	rt := &bearerTokenRoundTripper{rt: base, issuerHost: "issuer.example", tokenFilePath: tokenFile}
	rt.setJWKSHost("jwks.example:6443")

	cases := []struct {
		url  string
		want string
	}{
		{"https://issuer.example/x", "Bearer tok"},
		{"https://jwks.example:6443/x", "Bearer tok"},
		{"https://other.example/x", ""},
		{"https://jwks.example/x", ""}, // different port
		{"http://jwks.example:6443/x", ""},
	}
	for _, c := range cases {
		got = ""
		req, _ := http.NewRequest(http.MethodGet, c.url, nil)
		if _, err := rt.RoundTrip(req); err != nil {
			t.Fatal(err)
		}
		if got != c.want {
			t.Errorf("%s: Authorization=%q, want %q", c.url, got, c.want)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
