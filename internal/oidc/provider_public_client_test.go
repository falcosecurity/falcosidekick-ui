// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func exchangeAgainst(t *testing.T, secret string) (hdr http.Header, form map[string][]string) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		hdr = r.Header.Clone()
		form = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"at","token_type":"Bearer"}`))
	}))
	defer ts.Close()

	cfg := newOAuth2Config("cid", secret, "https://ui/cb", oauth2.Endpoint{AuthURL: ts.URL + "/auth", TokenURL: ts.URL + "/token"}, []string{"openid"})
	if _, err := cfg.Exchange(context.Background(), "code123", oauth2.VerifierOption("verif123")); err != nil {
		t.Fatalf("exchange: %v", err)
	}
	return hdr, form
}

func TestPublicClientTokenExchange(t *testing.T) {
	hdr, form := exchangeAgainst(t, "")
	if v := hdr.Get("Authorization"); v != "" {
		t.Errorf("unexpected Authorization header %q", v)
	}
	if form["client_id"][0] != "cid" {
		t.Errorf("client_id = %v", form["client_id"])
	}
	if form["code_verifier"][0] != "verif123" {
		t.Errorf("code_verifier = %v", form["code_verifier"])
	}
	if _, ok := form["client_secret"]; ok {
		t.Error("client_secret must not be sent")
	}
}

func TestConfidentialClientUnchanged(t *testing.T) {
	hdr, form := exchangeAgainst(t, "s3cret")
	if hdr.Get("Authorization") == "" {
		t.Error("expected Basic Authorization header for confidential client")
	}
	if _, ok := form["client_secret"]; ok {
		t.Error("secret should not be in body with default auto-detect")
	}
}
