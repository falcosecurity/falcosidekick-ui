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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/gomodule/redigo/redis"
	echo "github.com/labstack/echo/v4"
)

// OIDCTest constants for frequently used values in tests
const (
	//nolint:gosec // test fixture
	testOIDCRedirectURL      = "https://example.com/api/v1/auth/oidc/callback"
	testOIDCIssuer           = "https://example.com"
	testOIDCPathRules        = "/api/v1/rules"
	testKeyStatus            = "status"
	testGroupOIDC1           = "group1"
	testGroupOIDC2           = "group2"
	testUserAliceOIDC        = "alice"
	testEmailAliceOIDC       = "alice@example.com"
	testEmailBobOIDC         = "bob@example.com"
	testGroupAdminOIDC       = "admin"
	testGroupUsersOIDC       = "users"
	testGroupsAllowedOIDC    = "admin,operators"
	testXRequestedWithHeader = "XMLHttpRequest"
	testDefaultRedisAddr     = "localhost:6379"
)

// mockRedisConn is a simple mock redis.Conn for testing.
type mockRedisConn struct {
	data      map[string]string
	expiry    map[string]time.Time
	callCount int
}

func (m *mockRedisConn) Close() error {
	return nil
}

func (m *mockRedisConn) Err() error {
	return nil
}

func (m *mockRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	m.callCount++

	switch commandName {
	case "SETEX":
		if len(args) >= 3 {
			key := args[0].(string)
			ttl := args[1].(int)
			value := args[2].(string)
			m.data[key] = value
			m.expiry[key] = time.Now().Add(time.Duration(ttl) * time.Second)
			return "OK", nil
		}
		return nil, redis.ErrNil

	case methodGET:
		if len(args) >= 1 {
			key := args[0].(string)
			if val, ok := m.data[key]; ok {
				if exp, hasExp := m.expiry[key]; !hasExp || exp.After(time.Now()) {
					return []byte(val), nil
				}
				delete(m.data, key)
			}
		}
		return nil, redis.ErrNil

	case "GETDEL":
		if len(args) >= 1 {
			key := args[0].(string)
			if val, ok := m.data[key]; ok {
				delete(m.data, key)
				delete(m.expiry, key)
				return []byte(val), nil
			}
		}
		return nil, redis.ErrNil

	case "DEL":
		if len(args) >= 1 {
			key := args[0].(string)
			if _, ok := m.data[key]; ok {
				delete(m.data, key)
				delete(m.expiry, key)
				return int64(1), nil
			}
		}
		return int64(0), nil

	case "PTTL":
		if len(args) >= 1 {
			key := args[0].(string)
			if exp, ok := m.expiry[key]; ok {
				remaining := time.Until(exp).Milliseconds()
				if remaining > 0 {
					return remaining, nil
				}
				delete(m.data, key)
				delete(m.expiry, key)
			}
		}
		return nil, redis.ErrNil

	case "EXPIRE":
		if len(args) >= 2 {
			key := args[0].(string)
			ttl := args[1].(int)
			if _, ok := m.data[key]; ok {
				m.expiry[key] = time.Now().Add(time.Duration(ttl) * time.Second)
				return int64(1), nil
			}
		}
		return int64(0), nil

	default:
		return nil, redis.ErrNil
	}
}

func (m *mockRedisConn) Send(commandName string, args ...interface{}) error {
	_, _ = m.Do(commandName, args...)
	return nil
}

func (m *mockRedisConn) Flush() error {
	return nil
}

func (m *mockRedisConn) Receive() (interface{}, error) {
	return nil, redis.ErrNil
}

func newMockRedisConn() *mockRedisConn {
	return &mockRedisConn{
		data:   make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

// Test fixtures
func setupTestConfig(authMode, issuer, clientID string) {
	configuration.CreateConfiguration()
	config := configuration.GetConfiguration()
	config.AuthMode = authMode
	config.OIDCIssuer = issuer
	config.OIDCClientID = clientID
	config.OIDCRedirectURL = testOIDCRedirectURL
	config.OIDCUsernameClaim = claimPreferredUsername
	config.OIDCGroupsClaim = claimGroups
	config.SessionTTL = 28800        // 8 hours
	config.SessionIdleTimeout = 1800 // 30 minutes
}

// TestSessionMiddlewarePublicEndpoint tests that public endpoints skip session check
func TestSessionMiddlewarePublicEndpoint(t *testing.T) {
	setupTestConfig("oidc", testOIDCIssuer, "test-client")

	middleware := SessionMiddleware()

	nextHandler := func(c echo.Context) error {
		return nil
	}

	e := echo.New()

	tests := []struct {
		name   string
		path   string
		method string
		want   bool
	}{
		{name: "public healthz", path: pathHealthz, method: methodGET, want: true},
		{name: "public auth/me", path: pathAuthMe, method: methodGET, want: true},
		{name: "public login", path: pathLoginOIDC, method: methodGET, want: true},
		{name: "public callback", path: pathCallbackOIDC, method: methodGET, want: true},
		{name: "public ingestion POST /", path: pathIngestRoot, method: methodPOST, want: true},
		{name: "public ingestion POST /events/add", path: pathIngestAdd, method: methodPOST, want: true},
		{name: "public logout", path: pathLogout, method: methodPOST, want: true},
		{name: "protected GET /api/v1/rules", path: testOIDCPathRules, method: methodGET, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := middleware(nextHandler)(c)

			if err != nil && tt.want {
				t.Errorf("expected next to be called but got error: %v", err)
				return
			}

			if !tt.want && rec.Code == http.StatusOK {
				t.Errorf("protected endpoint %s should not allow access without session", tt.path)
			}
		})
	}
}

// TestSessionMiddlewareMissingSession tests that unauthenticated requests are rejected
func TestSessionMiddlewareMissingSession(t *testing.T) {
	setupTestConfig("oidc", testOIDCIssuer, "test-client")

	middleware := SessionMiddleware()

	nextCalled := false
	nextHandler := func(c echo.Context) error {
		nextCalled = true
		return c.JSON(http.StatusOK, map[string]string{testKeyStatus: "ok"})
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = middleware(nextHandler)(c)

	// Even though there's an error, the response should be 401
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}

	// Critical: next handler should NOT have been called
	if nextCalled {
		t.Errorf("SECURITY BUG: next handler was called after authentication failure")
	}

	// Verify response body
	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err == nil {
		if resp["error"] != "unauthenticated" {
			t.Errorf("expected error message 'unauthenticated', got %v", resp["error"])
		}
	}
}

// TestSessionMiddlewareNonOIDCMode tests that middleware is skipped when not in OIDC mode
func TestSessionMiddlewareNonOIDCMode(t *testing.T) {
	setupTestConfig("basic", testOIDCIssuer, "test-client")

	middleware := SessionMiddleware()

	nextCalled := false
	nextHandler := func(c echo.Context) error {
		nextCalled = true
		return c.JSON(http.StatusOK, map[string]string{testKeyStatus: "ok"})
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := middleware(nextHandler)(c)

	if err != nil {
		t.Errorf("expected no error in non-OIDC mode, got %v", err)
	}

	if !nextCalled {
		t.Errorf("next handler should be called in non-OIDC mode")
	}
}

// TestValidateConfig validates OIDC configuration
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid OIDC config",
			setup: func() {
				setupTestConfig("oidc", testOIDCIssuer, "test-client")
			},
			wantErr: false,
		},
		{
			name: "missing issuer",
			setup: func() {
				setupTestConfig("oidc", "", "test-client")
			},
			wantErr: true,
			errMsg:  "issuer is required",
		},
		{
			name: "http issuer without flag",
			setup: func() {
				setupTestConfig("oidc", "http://example.com", "test-client")
			},
			wantErr: true,
			errMsg:  "must be https",
		},
		{
			name: "missing client ID",
			setup: func() {
				configuration.CreateConfiguration()
				config := configuration.GetConfiguration()
				config.AuthMode = authModeOIDC
				config.OIDCIssuer = testOIDCIssuer
				config.OIDCClientID = ""
				config.OIDCRedirectURL = testOIDCRedirectURL
			},
			wantErr: true,
			errMsg:  "client ID is required",
		},
		{
			name: "missing redirect URL",
			setup: func() {
				configuration.CreateConfiguration()
				config := configuration.GetConfiguration()
				config.AuthMode = authModeOIDC
				config.OIDCIssuer = testOIDCIssuer
				config.OIDCClientID = "test-client"
				config.OIDCRedirectURL = ""
			},
			wantErr: true,
			errMsg:  "redirect URL is required",
		},
		{
			name: "wrong redirect URL path",
			setup: func() {
				configuration.CreateConfiguration()
				config := configuration.GetConfiguration()
				config.AuthMode = authModeOIDC
				config.OIDCIssuer = testOIDCIssuer
				config.OIDCClientID = "test-client"
				config.OIDCRedirectURL = "https://example.com/wrong/path"
			},
			wantErr: true,
			errMsg:  "must end with /api/v1/auth/oidc/callback",
		},
		{
			name: "non-OIDC mode skips validation",
			setup: func() {
				configuration.CreateConfiguration()
				config := configuration.GetConfiguration()
				config.AuthMode = "basic"
				config.OIDCIssuer = ""
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := ValidateConfig()

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateConfig() error message %q, want to contain %q", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestExtractClaims tests claim extraction with various claim structures
func TestExtractClaims(t *testing.T) {
	setupTestConfig("oidc", testOIDCIssuer, "test-client")

	tests := []struct {
		name           string
		claims         map[string]interface{}
		expectUsername string
		expectEmail    string
		expectGroups   []string
	}{
		{
			name: "preferred_username claim",
			claims: map[string]interface{}{
				claimPreferredUsername: testUserAliceOIDC,
				"email":                testEmailAliceOIDC,
				claimGroups:            []interface{}{testGroupOIDC1, testGroupOIDC2},
			},
			expectUsername: testUserAliceOIDC,
			expectEmail:    testEmailAliceOIDC,
			expectGroups:   []string{testGroupOIDC1, testGroupOIDC2},
		},
		{
			name: "fallback to email when preferred_username missing",
			claims: map[string]interface{}{
				"email":     testEmailBobOIDC,
				claimGroups: []interface{}{testGroupOIDC1},
			},
			expectUsername: testEmailBobOIDC,
			expectEmail:    testEmailBobOIDC,
			expectGroups:   []string{testGroupOIDC1},
		},
		{
			name: "fallback to sub when username and email missing",
			claims: map[string]interface{}{
				claimSub: ingestTestUserID,
			},
			expectUsername: ingestTestUserID,
			expectEmail:    "",
			expectGroups:   []string{},
		},
		{
			name: "groups as string instead of array",
			claims: map[string]interface{}{
				claimPreferredUsername: "charlie",
				claimGroups:            testGroupAdminOIDC,
			},
			expectUsername: "charlie",
			expectGroups:   []string{testGroupAdminOIDC},
		},
		{
			name: "groups as array with non-string values (filtered)",
			claims: map[string]interface{}{
				claimPreferredUsername: "dave",
				claimGroups:            []interface{}{testGroupOIDC1, 123, testGroupOIDC2},
			},
			expectUsername: "dave",
			expectGroups:   []string{testGroupOIDC1, testGroupOIDC2},
		},
		{
			name: "empty claims",
			claims: map[string]interface{}{
				claimSub: "",
			},
			expectUsername: "",
			expectEmail:    "",
			expectGroups:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockToken := &mockIDToken{claims: tt.claims}

			username, email, groups := extractClaims(mockToken)

			if username != tt.expectUsername {
				t.Errorf("username = %q, want %q", username, tt.expectUsername)
			}

			if email != tt.expectEmail {
				t.Errorf("email = %q, want %q", email, tt.expectEmail)
			}

			if len(groups) != len(tt.expectGroups) {
				t.Errorf("groups length = %d, want %d", len(groups), len(tt.expectGroups))
			}

			for i, g := range groups {
				if i < len(tt.expectGroups) && g != tt.expectGroups[i] {
					t.Errorf("groups[%d] = %q, want %q", i, g, tt.expectGroups[i])
				}
			}
		})
	}
}

// TestIsUserInAllowedGroups tests group authorization logic
func TestIsUserInAllowedGroups(t *testing.T) {
	tests := []struct {
		name       string
		userGroups []string
		allowedStr string
		expect     bool
	}{
		{
			name:       "user in allowed group",
			userGroups: []string{testGroupAdminOIDC, testGroupUsersOIDC},
			allowedStr: testGroupsAllowedOIDC,
			expect:     true,
		},
		{
			name:       "user not in allowed groups",
			userGroups: []string{testGroupUsersOIDC},
			allowedStr: testGroupsAllowedOIDC,
			expect:     false,
		},
		{
			name:       "empty user groups",
			userGroups: []string{},
			allowedStr: testGroupAdminOIDC,
			expect:     false,
		},
		{
			name:       "groups with whitespace",
			userGroups: []string{testGroupAdminOIDC},
			allowedStr: "  admin  , operators ",
			expect:     true,
		},
		{
			name:       "multiple matching groups",
			userGroups: []string{testGroupAdminOIDC, "operators", testGroupUsersOIDC},
			allowedStr: testGroupsAllowedOIDC,
			expect:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUserInAllowedGroups(tt.userGroups, tt.allowedStr)
			if result != tt.expect {
				t.Errorf("isUserInAllowedGroups() = %v, want %v", result, tt.expect)
			}
		})
	}
}

// TestGetAppBase tests URL base extraction
func TestGetAppBase(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "standard redirect URL",
			input:  "https://example.com/api/v1/auth/oidc/callback",
			expect: testOIDCIssuer,
		},
		{
			name:   "with subpath",
			input:  "https://example.com/app/api/v1/auth/oidc/callback",
			expect: "https://example.com/app",
		},
		{
			name:   "missing callback path",
			input:  "https://example.com/some/path",
			expect: "https://example.com/some/path",
		},
		{
			name:   "localhost",
			input:  "http://localhost:3000/api/v1/auth/oidc/callback",
			expect: "http://localhost:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAppBase(tt.input)
			if result != tt.expect {
				t.Errorf("getAppBase(%q) = %q, want %q", tt.input, result, tt.expect)
			}
		})
	}
}

// TestExtractOrigin tests origin extraction from URLs
func TestExtractOrigin(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "https URL",
			input:  "https://example.com/api/v1/auth/oidc/callback",
			expect: testOIDCIssuer,
		},
		{
			name:   "http URL",
			input:  "http://localhost:3000/some/path",
			expect: "http://localhost:3000",
		},
		{
			name:   "URL with port",
			input:  "https://example.com:8443/path",
			expect: "https://example.com:8443",
		},
		{
			name:   "IP address",
			input:  "https://192.168.1.1/path",
			expect: "https://192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractOrigin(tt.input)
			if result != tt.expect {
				t.Errorf("extractOrigin(%q) = %q, want %q", tt.input, result, tt.expect)
			}
		})
	}
}

// TestIsPublicEndpoint tests endpoint visibility
func TestIsPublicEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		method string
		expect bool
	}{
		{name: "GET /api/v1/healthz", path: pathHealthz, method: methodGET, expect: true},
		{name: "GET /api/v1/auth/me", path: pathAuthMe, method: methodGET, expect: true},
		{name: "GET /api/v1/auth/oidc/login", path: pathLoginOIDC, method: methodGET, expect: true},
		{name: "GET /api/v1/auth/oidc/callback", path: pathCallbackOIDC, method: methodGET, expect: true},
		{name: "POST /api/v1/", path: pathIngestRoot, method: methodPOST, expect: true},
		{name: "POST /api/v1/events/add", path: pathIngestAdd, method: methodPOST, expect: true},
		{name: "POST /api/v1/auth/logout", path: pathLogout, method: methodPOST, expect: true},
		{name: "GET /api/v1/rules", path: testOIDCPathRules, method: methodGET, expect: false},
		{name: "POST /api/v1/rules", path: testOIDCPathRules, method: methodPOST, expect: false},
		{name: "DELETE /api/v1/rules/1", path: "/api/v1/rules/1", method: "DELETE", expect: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPublicEndpoint(tt.path, tt.method)
			if result != tt.expect {
				t.Errorf("isPublicEndpoint(%q, %q) = %v, want %v", tt.path, tt.method, result, tt.expect)
			}
		})
	}
}

// TestSessionStorageAndRetrieval tests session CRUD operations with mock Redis
func TestSessionStorageAndRetrieval(t *testing.T) {
	conn := newMockRedisConn()

	// Generate session ID
	sessionID, err := GenerateID()
	if err != nil {
		t.Fatalf("failed to generate session ID: %v", err)
	}

	// Create session
	session := &Session{
		Sub:       ingestTestUserID,
		Username:  testUserAliceOIDC,
		Email:     testEmailAliceOIDC,
		Groups:    []string{testGroupAdminOIDC, testGroupUsersOIDC},
		IDToken:   "eyJ...",
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}

	// Store session
	err = StoreSession(conn, sessionID, session)
	if err != nil {
		t.Fatalf("failed to store session: %v", err)
	}

	// Retrieve session
	retrieved, err := GetSession(conn, sessionID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	if retrieved.Username != session.Username {
		t.Errorf("username = %q, want %q", retrieved.Username, session.Username)
	}

	if len(retrieved.Groups) != len(session.Groups) {
		t.Errorf("groups count = %d, want %d", len(retrieved.Groups), len(session.Groups))
	}

	// Delete session
	err = DeleteSession(conn, sessionID)
	if err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	// Verify deleted
	_, err = GetSession(conn, sessionID)
	if err == nil {
		t.Errorf("expected error after deletion, got none")
	}
}

// TestFlowStorageAndRetrieval tests flow CRUD operations
func TestFlowStorageAndRetrieval(t *testing.T) {
	conn := newMockRedisConn()

	flowID, err := GenerateID()
	if err != nil {
		t.Fatalf("failed to generate flow ID: %v", err)
	}

	flow := &Flow{
		State:     "state123",
		Nonce:     "nonce456",
		Verifier:  "verifier789",
		CreatedAt: time.Now(),
	}

	// Store flow
	err = StoreFlow(conn, flowID, flow)
	if err != nil {
		t.Fatalf("failed to store flow: %v", err)
	}

	// Retrieve flow (GETDEL - one-time use)
	retrieved, err := GetFlow(conn, flowID)
	if err != nil {
		t.Fatalf("failed to get flow: %v", err)
	}

	if retrieved.State != flow.State {
		t.Errorf("state = %q, want %q", retrieved.State, flow.State)
	}

	// Verify flow is deleted (one-time use)
	_, err = GetFlow(conn, flowID)
	if err == nil {
		t.Errorf("expected error on second retrieval (GETDEL), got none")
	}
}

// TestHashID verifies consistent hashing
func TestHashID(t *testing.T) {
	id := "test-session-id-12345"

	hash1 := HashID(id)
	hash2 := HashID(id)

	if hash1 != hash2 {
		t.Errorf("HashID not consistent: %q != %q", hash1, hash2)
	}

	if hash1 == id {
		t.Errorf("HashID should not return plaintext ID")
	}

	if len(hash1) == 0 {
		t.Errorf("HashID returned empty string")
	}
}

// Mock ID token for testing
type mockIDToken struct {
	claims map[string]interface{}
}

func (m *mockIDToken) Claims(v interface{}) error {
	data, _ := json.Marshal(m.claims)
	return json.Unmarshal(data, v)
}

// Helper function
func contains(s, substr string) bool {
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestLogoutOriginValidation tests that logout properly validates Origin header and X-Requested-With.
func TestLogoutOriginValidation(t *testing.T) {
	setupTestConfig("oidc", testOIDCIssuer, "test-client")
	config := configuration.GetConfiguration()
	config.OIDCRedirectURL = testOIDCRedirectURL

	e := echo.New()

	tests := []struct {
		name       string
		origin     string
		xreqwith   string
		expectCode int
	}{
		{
			name:       "valid origin",
			origin:     "https://example.com",
			xreqwith:   testXRequestedWithHeader,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "origin case-insensitive match",
			origin:     "https://UI.example.com",
			xreqwith:   testXRequestedWithHeader,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "evil origin",
			origin:     "https://evil.com",
			xreqwith:   testXRequestedWithHeader,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "no origin but X-Requested-With present",
			origin:     "",
			xreqwith:   testXRequestedWithHeader,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "no X-Requested-With",
			origin:     "",
			xreqwith:   "",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.xreqwith != "" {
				req.Header.Set("X-Requested-With", tt.xreqwith)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = Logout(c)
		})
	}
}

// TestRefreshSessionTTLExpired tests that RefreshSessionTTL returns ErrSessionExpired when session has exceeded idle timeout.
func TestRefreshSessionTTLExpired(t *testing.T) {
	setupTestConfig("oidc", testOIDCIssuer, "test-client")

	conn := newMockRedisConn()

	// Create a session that has exceeded idle timeout
	sessionID := "test-session"
	session := &Session{
		Username:  "testuser",
		CreatedAt: time.Now().Add(-1 * time.Hour),
		LastSeen:  time.Now().Add(-2 * time.Hour), // 2 hours ago, exceeds 30-min idle timeout
	}

	// Store the session in mock Redis
	data, _ := json.Marshal(session)
	key := fmt.Sprintf("fsui:session:%s", HashID(sessionID))
	conn.data[key] = string(data)
	conn.expiry[key] = time.Now().Add(time.Hour)

	// Set idle timeout to 30 minutes
	config := configuration.GetConfiguration()
	config.SessionIdleTimeout = 1800 // 30 minutes

	// RefreshSessionTTL should return ErrSessionExpired
	err := RefreshSessionTTL(conn, sessionID)
	if err == nil || !errors.Is(err, ErrSessionExpired) {
		t.Errorf("RefreshSessionTTL() expected ErrSessionExpired, got %v", err)
	}
}

// TestRefreshSessionTTLTransientError tests that transient Redis errors don't delete the session.
func TestRefreshSessionTTLTransientError(t *testing.T) {
	setupTestConfig("oidc", testOIDCIssuer, "test-client")

	// Create a mock Redis connection that returns an error on SET command
	mockConn := &mockRedisConnWithError{
		data:   make(map[string]string),
		expiry: make(map[string]time.Time),
		errOn:  "SET",
	}

	// Create a valid session
	sessionID := "test-session"
	session := &Session{
		Username:  "testuser",
		CreatedAt: time.Now(),
		LastSeen:  time.Now().Add(-70 * time.Second), // 70 seconds ago - triggers refresh threshold
	}

	// Store the session
	data, _ := json.Marshal(session)
	key := fmt.Sprintf("fsui:session:%s", HashID(sessionID))
	mockConn.data[key] = string(data)
	mockConn.expiry[key] = time.Now().Add(time.Hour)

	// Set idle timeout
	config := configuration.GetConfiguration()
	config.SessionIdleTimeout = 1800 // 30 minutes

	// RefreshSessionTTL should return a transient error (not ErrSessionExpired)
	err := RefreshSessionTTL(mockConn, sessionID)
	if err == nil || errors.Is(err, ErrSessionExpired) {
		t.Errorf("RefreshSessionTTL() expected transient error, got %v", err)
	}

	// Session should still exist (not deleted)
	if _, ok := mockConn.data[key]; !ok {
		t.Errorf("session was deleted on transient error")
	}
}

// mockRedisConnWithError is a mock Redis connection that injects errors on specific commands.
type mockRedisConnWithError struct {
	data      map[string]string
	expiry    map[string]time.Time
	callCount int
	errOn     string
}

func (m *mockRedisConnWithError) Close() error {
	return nil
}

func (m *mockRedisConnWithError) Err() error {
	return nil
}

func (m *mockRedisConnWithError) Do(commandName string, args ...interface{}) (interface{}, error) {
	m.callCount++

	if commandName == m.errOn {
		return nil, fmt.Errorf("injected error on %s", commandName)
	}

	switch commandName {
	case "SETEX":
		if len(args) >= 3 {
			key := args[0].(string)
			ttl := args[1].(int)
			value := args[2].(string)
			m.data[key] = value
			m.expiry[key] = time.Now().Add(time.Duration(ttl) * time.Second)
			return "OK", nil
		}
		return nil, redis.ErrNil

	case methodGET:
		if len(args) >= 1 {
			key := args[0].(string)
			if val, ok := m.data[key]; ok {
				if exp, hasExp := m.expiry[key]; !hasExp || exp.After(time.Now()) {
					return []byte(val), nil
				}
				delete(m.data, key)
			}
		}
		return nil, redis.ErrNil

	case "SET":
		if len(args) >= 2 {
			key := args[0].(string)
			value := args[1].(string)
			// Handle SET key value XX EX ttl
			if len(args) >= 4 && args[2] == "XX" && args[3] == "EX" {
				if _, ok := m.data[key]; ok {
					if len(args) >= 5 {
						ttl := args[4].(int)
						m.data[key] = value
						m.expiry[key] = time.Now().Add(time.Duration(ttl) * time.Second)
						return "OK", nil
					}
				}
				return nil, redis.ErrNil
			}
			m.data[key] = value
			return "OK", nil
		}
		return nil, redis.ErrNil

	default:
		return nil, redis.ErrNil
	}
}

func (m *mockRedisConnWithError) Send(commandName string, args ...interface{}) error {
	_, _ = m.Do(commandName, args...)
	return nil
}

func (m *mockRedisConnWithError) Flush() error {
	return nil
}

func (m *mockRedisConnWithError) Receive() (interface{}, error) {
	return nil, redis.ErrNil
}

// TestResolveRedisAddr tests the resolveRedisAddr helper function.
func TestResolveRedisAddr(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "hostname only", input: "redis-master", expect: "redis-master:6379"},
		{name: "empty host with port", input: ":6380", expect: "localhost:6380"},
		{name: "host and port", input: "h:1", expect: "h:1"},
		{name: "empty string", input: "", expect: testDefaultRedisAddr},
		{name: "hostname with default port", input: "localhost:6379", expect: testDefaultRedisAddr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveRedisAddr(tt.input)
			if result != tt.expect {
				t.Errorf("resolveRedisAddr(%q) = %q, want %q", tt.input, result, tt.expect)
			}
		})
	}
}
