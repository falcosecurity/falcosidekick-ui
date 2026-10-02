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
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// Auth mode constants
const (
	authModeOIDC  = "oidc"
	routeNameAuth = "authenticate"
)

// RegisterRoutes registers OIDC-mode auth routes and middleware in the given echo group.
// This includes SessionMiddleware, the auth endpoints (Me, Login, Callback, Logout),
// and 404 handlers for /auth and /authenticate.
func RegisterRoutes(g *echo.Group) {
	g.Use(SessionMiddleware())

	// OIDC auth routes
	g.GET("/auth/me", Me).Name = "auth-me"
	g.GET("/auth/oidc/login", Login).Name = "oidc-login"
	g.GET("/auth/oidc/callback", Callback).Name = "oidc-callback"
	g.POST("/auth/logout", Logout).Name = "auth-logout"

	// In OIDC mode, basic auth endpoints return 404
	g.POST("/auth", func(c echo.Context) error {
		noStore(c)
		return echo.NewHTTPError(http.StatusNotFound)
	}).Name = routeNameAuth
	g.POST("/authenticate", func(c echo.Context) error {
		noStore(c)
		return echo.NewHTTPError(http.StatusNotFound)
	}).Name = routeNameAuth
}

// Response header constants
const (
	ErrorJSON               = "error"
	headerXRequestedWith    = "X-Requested-With"
	headerXRequestedWithVal = "XMLHttpRequest"
)

// OIDC claim constants
const (
	claimPreferredUsername = "preferred_username"
	claimEmail             = "email"
	claimSub               = "sub"
	claimGroups            = "groups"
	claimAzp               = "azp"
	claimClientID          = "client_id"
	claimAud               = "aud"
	claimExp               = "exp"
)

// OIDC error message constants
const (
	errorUnauthenticated = "unauthenticated"
	errorUnauthorized    = "unauthorized"
)

// HTTP scheme constants
//
//nolint:goconst
const (
	schemeHTTPS = "https"
	schemeHTTP  = "http"
)

// setCookie sets a cookie with consistent security attributes.
// G124 is suppressed because insecure mode is intentional for development.
//
//nolint:gosec
func setCookie(c echo.Context, name string, value string, maxAge int, config *configuration.Configuration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   !config.OIDCInsecureAllowHTTP,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
	c.SetCookie(cookie)
}

// Login generates OIDC flow and redirects to provider.
// @Summary      OIDC Login
// @Description  Start OIDC login flow
// @Produce      json
// @Success      302  "Redirect to IdP"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /api/v1/auth/oidc/login [get]
func Login(c echo.Context) error {
	ctx := c.Request().Context()
	config := configuration.GetConfiguration()

	// Generate flow components
	flowID, err := GenerateID()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate flow ID")
	}

	state, err := GenerateID()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate state")
	}

	nonce, err := GenerateID()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate nonce")
	}

	verifier := oauth2.GenerateVerifier()

	// Store flow in Redis
	flow := &Flow{
		State:     state,
		Nonce:     nonce,
		Verifier:  verifier,
		CreatedAt: time.Now(),
	}

	conn, err := GetRedisConn()
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to get redis connection: %v", err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to store flow")
	}
	defer conn.Close()

	err = StoreFlow(conn, flowID, flow)
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to store flow: %v", err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to store flow")
	}

	// Get OAuth2 config
	oauth2Cfg, err := GetOAuth2Config(ctx)
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to get OAuth2 config: %v", err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize OAuth2")
	}

	// Generate authorization URL (use state, not flowID)
	authURL := oauth2Cfg.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
		oidc.Nonce(nonce),
	)

	// Set flow cookie (use flowID, not state)
	flowCookieName := "__Host-fsui_oidc_flow"
	if config.OIDCInsecureAllowHTTP {
		flowCookieName = "fsui_oidc_flow"
	}
	setCookie(c, flowCookieName, flowID, 600, config)

	// Set no-store cache headers
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Pragma", "no-cache")

	return c.Redirect(http.StatusFound, authURL)
}

// Callback handles the OAuth2 callback from the IdP.
// @Summary      OIDC Callback
// @Description  Handle OIDC callback from IdP
// @Param        code    query  string  false  "Authorization code"
// @Param        state   query  string  false  "State parameter"
// @Param        error   query  string  false  "Error from IdP"
// @Param        iss     query  string  false  "Issuer"
// @Produce      json
// @Success      302  "Redirect to dashboard"
// @Failure      400  {string}  string  "Bad Request"
// @Failure      403  {string}  string  "Forbidden"
// @Router       /api/v1/auth/oidc/callback [get]
func Callback(c echo.Context) error {
	ctx := c.Request().Context()
	config := configuration.GetConfiguration()

	// Get flow cookie
	cookie, err := c.Cookie("__Host-fsui_oidc_flow")
	if err != nil && config.OIDCInsecureAllowHTTP {
		cookie, err = c.Cookie("fsui_oidc_flow")
	}

	// Determine flow cookie name for clearing
	flowCookieName := "__Host-fsui_oidc_flow"
	if config.OIDCInsecureAllowHTTP {
		flowCookieName = "fsui_oidc_flow"
	}

	// Clear flow cookie at the very START if present (before GetFlow)
	if err == nil && cookie != nil {
		setCookie(c, flowCookieName, "", -1, config)
	}

	if err != nil {
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Get flow from Redis
	conn, err := GetRedisConn()
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to get redis connection: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}
	defer conn.Close()

	flow, err := GetFlow(conn, cookie.Value)
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to get flow: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Check for error from IdP
	// Validate error parameter format before logging
	if idpErr := c.QueryParam("error"); idpErr != "" {
		if isValidOIDCErrorCode(idpErr) {
			utils.WriteLog("error", fmt.Sprintf("IdP error: %s", idpErr))
		} else {
			utils.WriteLog("error", "IdP error: invalid")
		}
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Verify RFC 9207: if iss parameter is present, it must match issuer
	if iss := c.QueryParam("iss"); iss != "" {
		if subtle.ConstantTimeCompare([]byte(iss), []byte(config.OIDCIssuer)) != 1 {
			c.Response().Header().Set("Cache-Control", "no-store")
			c.Response().Header().Set("Pragma", "no-cache")
			return echo.NewHTTPError(http.StatusBadRequest, "issuer mismatch")
		}
	}

	// Verify state
	state := c.QueryParam("state")
	if subtle.ConstantTimeCompare([]byte(state), []byte(flow.State)) != 1 {
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return echo.NewHTTPError(http.StatusBadRequest, "state mismatch")
	}

	// Exchange code for token
	code := c.QueryParam("code")
	oauth2Cfg, err := GetOAuth2Config(ctx)
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to get OAuth2 config: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Create context with custom HTTP client for token exchange
	httpClient, err := getHTTPClient()
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to create HTTP client: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	token, err := oauth2Cfg.Exchange(ctx, code, oauth2.VerifierOption(flow.Verifier))
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to exchange token: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Extract ID token
	rawIDToken := token.Extra("id_token")
	if rawIDToken == nil {
		utils.WriteLog("error", "missing id_token in token response")
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	idTokenString, ok := rawIDToken.(string)
	if !ok {
		utils.WriteLog("error", "id_token is not a string")
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Verify ID token
	verifier, err := GetIDTokenVerifier(ctx)
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to get ID token verifier: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	idToken, err := verifier.Verify(ctx, idTokenString)
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to verify ID token: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Verify nonce
	var claims struct {
		Nonce string `json:"nonce"`
		AZP   string `json:"azp"`
	}
	if err := idToken.Claims(&claims); err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to parse claims: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	if subtle.ConstantTimeCompare([]byte(claims.Nonce), []byte(flow.Nonce)) != 1 {
		utils.WriteLog("error", "nonce mismatch")
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Verify azp if aud has multiple entries
	var allClaims map[string]interface{}
	if err := idToken.Claims(&allClaims); err == nil {
		aud, ok := allClaims["aud"].([]interface{})
		if ok && len(aud) > 1 && claims.AZP != "" {
			if subtle.ConstantTimeCompare([]byte(claims.AZP), []byte(config.OIDCClientID)) != 1 {
				utils.WriteLog("error", "azp mismatch")
				c.Response().Header().Set("Cache-Control", "no-store")
				c.Response().Header().Set("Pragma", "no-cache")
				return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
			}
		}
	}

	// Extract username, email, and groups
	username, email, groups := extractClaims(idToken)

	// Check allowed groups
	if config.OIDCAllowedGroups != "" {
		if !isUserInAllowedGroups(groups, config.OIDCAllowedGroups) {
			utils.WriteLog("warning", fmt.Sprintf("user %v not in allowed groups", username))
			c.Response().Header().Set("Cache-Control", "no-store")
			c.Response().Header().Set("Pragma", "no-cache")
			return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=forbidden")
		}
	}

	utils.WriteLog("info", fmt.Sprintf("user '%v' authenticated via OIDC", username))

	// Delete existing session if present (session fixation prevention)
	existingSessionCookie, err := c.Cookie("__Host-fsui_session")
	if err != nil && config.OIDCInsecureAllowHTTP {
		existingSessionCookie, err = c.Cookie("fsui_session")
	}
	if err == nil && existingSessionCookie != nil {
		DeleteSession(conn, existingSessionCookie.Value)
	}

	// Create new session
	sessionID, err := GenerateID()
	if err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to generate session ID: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	session := &Session{
		Sub:       idToken.Subject,
		Username:  username,
		Email:     email,
		Groups:    groups,
		IDToken:   idTokenString,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}

	if err := StoreSession(conn, sessionID, session); err != nil {
		utils.WriteLog("error", fmt.Sprintf("failed to store session: %v", err))
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")
		return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL)+"#/login?error=sso_failed")
	}

	// Set session cookie
	sessionCookieName := "__Host-fsui_session"
	if config.OIDCInsecureAllowHTTP {
		sessionCookieName = "fsui_session"
	}
	setCookie(c, sessionCookieName, sessionID, 0, config)

	// Redirect to app base
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Pragma", "no-cache")
	return c.Redirect(http.StatusFound, getAppBase(config.OIDCRedirectURL))
}

// Logout handles OIDC logout.
// @Summary      OIDC Logout
// @Description  Logout from OIDC
// @Produce      json
// @Success      200  {object}  map[string]string  "Logout response"
// @Failure      400  {string}  string             "Bad Request"
// @Failure      401  {string}  string             "Unauthorized"
// @Router       /api/v1/auth/logout [post]
func Logout(c echo.Context) error {
	config := configuration.GetConfiguration()

	// CSRF defense
	if c.Request().Header.Get(headerXRequestedWith) != headerXRequestedWithVal {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	// Verify origin if present, using X-Requested-With as fallback
	if origin := c.Request().Header.Get("Origin"); origin != "" {
		redirectOrigin := extractOrigin(config.OIDCRedirectURL)
		// Normalize origin for comparison (case-insensitive host, default ports)
		normalizedOrigin := extractOrigin(origin)
		if subtle.ConstantTimeCompare([]byte(normalizedOrigin), []byte(redirectOrigin)) != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "origin mismatch")
		}
	}

	// Get session cookie
	sessionCookie, err := c.Cookie("__Host-fsui_session")
	if err != nil && config.OIDCInsecureAllowHTTP {
		sessionCookie, err = c.Cookie("fsui_session")
	}

	// Read session to extract id_token BEFORE deletion
	var idToken string
	if err == nil {
		conn, err := GetRedisConn()
		if err == nil {
			defer conn.Close()
			if session, err := GetSession(conn, sessionCookie.Value); err == nil && session.IDToken != "" {
				idToken = session.IDToken
			}
		}
	}

	// Delete session from Redis
	if err == nil {
		conn, err := GetRedisConn()
		if err == nil {
			defer conn.Close()
			DeleteSession(conn, sessionCookie.Value)
		}
	}

	// Compute cookie name for expiration
	cookieName := "__Host-fsui_session"
	if config.OIDCInsecureAllowHTTP {
		cookieName = "fsui_session"
	}

	// Use setCookie helper to respect INSECURE_ALLOW_HTTP
	setCookie(c, cookieName, "", -1, config)

	// Get end_session_endpoint from provider
	logoutURL := ""
	if config.OIDCPostLogoutRedirect != "" {
		ctx := c.Request().Context()
		provider, err := GetProvider(ctx)
		if err == nil {
			var claims struct {
				EndSessionEndpoint string `json:"end_session_endpoint"`
			}
			if provider.Claims(&claims) == nil && claims.EndSessionEndpoint != "" {
				params := url.Values{}
				params.Set("post_logout_redirect_uri", config.OIDCPostLogoutRedirect)
				params.Set("client_id", config.OIDCClientID)
				if idToken != "" {
					params.Set("id_token_hint", idToken)
				}
				logoutURL = claims.EndSessionEndpoint + "?" + params.Encode()
			}
		}
	}

	// Set no-store cache headers
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Pragma", "no-cache")

	return c.JSON(http.StatusOK, map[string]string{"logout_url": logoutURL})
}

// Me returns authentication information.
// @Summary      Get Auth Info
// @Description  Get current authentication mode and status
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Auth info"
// @Router       /api/v1/auth/me [get]
func Me(c echo.Context) error {
	config := configuration.GetConfiguration()

	response := map[string]interface{}{
		"mode":          config.AuthMode,
		"authenticated": false,
		"username":      "",
	}

	if config.AuthMode == authModeOIDC {
		// Check for session cookie
		sessionCookie, err := c.Cookie("__Host-fsui_session")
		if err != nil && config.OIDCInsecureAllowHTTP {
			sessionCookie, err = c.Cookie("fsui_session")
		}

		if err == nil {
			conn, err := GetRedisConn()
			if err == nil {
				defer conn.Close()

				session, err := GetSession(conn, sessionCookie.Value)
				if err == nil {
					// Check absolute TTL
					createdAt := session.CreatedAt
					ttl := time.Duration(config.SessionTTL) * time.Second

					if time.Since(createdAt) < ttl {
						response["authenticated"] = true
						response["username"] = session.Username
					}
				}
			}
		}
	}

	// Set no-store cache headers
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Pragma", "no-cache")

	return c.JSON(http.StatusOK, response)
}

// Helper functions

// noStore sets Cache-Control and Pragma headers to prevent caching.
func noStore(c echo.Context) {
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Pragma", "no-cache")
}

func extractClaims(idToken interface{ Claims(v interface{}) error }) (username, email string, groups []string) {
	config := configuration.GetConfiguration()

	var allClaims map[string]interface{}
	if err := idToken.Claims(&allClaims); err != nil {
		return "", "", []string{}
	}

	// Extract username
	usernameClaim := config.OIDCUsernameClaim
	if usernameClaim == "" {
		usernameClaim = claimPreferredUsername
	}

	if u, ok := allClaims[usernameClaim].(string); ok && u != "" {
		username = u
	} else if e, ok := allClaims[claimEmail].(string); ok && e != "" {
		username = e
	} else if s, ok := allClaims[claimSub].(string); ok && s != "" {
		username = s
	}

	// Extract email
	if e, ok := allClaims[claimEmail].(string); ok {
		email = e
	}

	// Extract groups
	groupsClaim := config.OIDCGroupsClaim
	if groupsClaim == "" {
		groupsClaim = claimGroups
	}

	if g, ok := allClaims[groupsClaim].([]interface{}); ok {
		for _, v := range g {
			if s, ok := v.(string); ok {
				groups = append(groups, s)
			}
		}
	} else if g, ok := allClaims[groupsClaim].(string); ok {
		groups = []string{g}
	}

	return username, email, groups
}

func isUserInAllowedGroups(userGroups []string, allowedGroupsStr string) bool {
	// Reject empty group values in user's groups
	for _, userGroup := range userGroups {
		if userGroup == "" {
			continue
		}
		// Check if this group is in allowlist
		allowedGroups := strings.Split(allowedGroupsStr, ",")
		for _, allowedGroup := range allowedGroups {
			trimmed := strings.TrimSpace(allowedGroup)
			// Skip empty entries after trim
			if trimmed == "" {
				continue
			}
			if userGroup == trimmed {
				return true
			}
		}
	}

	return false
}

func getAppBase(redirectURL string) string {
	// Remove /api/v1/auth/oidc/callback from redirect URL
	parts := strings.Split(redirectURL, "/api/v1/auth/oidc/callback")
	if len(parts) > 0 {
		return parts[0]
	}
	return redirectURL
}

func extractOrigin(urlStr string) string {
	// Parse URL and normalize with default ports
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	host := u.Hostname()
	port := u.Port()

	// Normalize ports: omit default ports
	if port == "" {
		if u.Scheme == schemeHTTPS {
			port = "443"
		} else if u.Scheme == schemeHTTP {
			port = "80"
		}
	}

	// Only include port if it's non-default
	isDefaultPort := (u.Scheme == schemeHTTPS && port == "443") ||
		(u.Scheme == schemeHTTP && port == "80")

	if isDefaultPort {
		return u.Scheme + "://" + strings.ToLower(host)
	}

	return u.Scheme + "://" + strings.ToLower(host) + ":" + port
}

// isValidOIDCErrorCode validates error code format per OAuth 2.0 / OpenID Connect specs.
// Only log error codes matching ^[a-z_]{1,64}$, else log "invalid"
func isValidOIDCErrorCode(code string) bool {
	if len(code) == 0 || len(code) > 64 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-z_]{1,64}$`, code)
	return matched
}
