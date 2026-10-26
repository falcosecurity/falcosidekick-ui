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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"golang.org/x/oauth2"
)

// bearerTokenRoundTripper wraps an http.RoundTripper and adds Authorization header
// only for requests to the issuer host and to the host of the jwks_uri advertised by
// discovery (they differ e.g. on Kubernetes, where the issuer is
// https://kubernetes.default.svc.cluster.local but jwks_uri points at the API server
// address). It never sends the token to any other host. It re-reads the token file on each request
// to support token rotation (e.g., Kubernetes projected tokens).
type bearerTokenRoundTripper struct {
	rt                http.RoundTripper
	issuerHost        string
	jwksHost          string // set after discovery via setJWKSHost; guarded by mu
	tokenFilePath     string
	lastModTime       time.Time
	lastCheck         time.Time
	cachedToken       string
	mu                sync.Mutex
	insecureAllowHTTP bool
}

func (b *bearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request before modifying headers
	req = req.Clone(req.Context())

	// Add Authorization header if request is to issuer host
	// Only allow HTTPS by default, or HTTP if InsecureAllowHTTP is set
	//nolint:goconst
	isSecure := req.URL.Scheme == "https" || (b.insecureAllowHTTP && req.URL.Scheme == "http")
	if isSecure && b.isTrustedHost(req.URL.Host) {
		token := b.getToken()
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
	return b.rt.RoundTrip(req)
}

// isTrustedHost reports whether the bearer token may be sent to host.
func (b *bearerTokenRoundTripper) isTrustedHost(host string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return host == b.issuerHost || (b.jwksHost != "" && host == b.jwksHost)
}

// setJWKSHost records the host of the discovered jwks_uri as a second trusted host.
func (b *bearerTokenRoundTripper) setJWKSHost(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.jwksHost = host
}

// getToken re-reads the bearer token file at most every 30s to support rotation
func (b *bearerTokenRoundTripper) getToken() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()

	// Only stat the file at most every 30s based on lastCheck
	if now.Sub(b.lastCheck) < 30*time.Second {
		return b.cachedToken
	}

	// Update lastCheck timestamp
	b.lastCheck = now

	// Try to read and stat the file
	stat, err := os.Stat(b.tokenFilePath)
	if err != nil {
		// If we can't stat, keep using cached token
		return b.cachedToken
	}

	// If file hasn't been modified and we have a cached token, use it
	if stat.ModTime() == b.lastModTime && b.cachedToken != "" {
		return b.cachedToken
	}

	// Read the file
	data, err := os.ReadFile(b.tokenFilePath)
	if err != nil {
		// If we can't read, keep using cached token
		return b.cachedToken
	}

	// Trim whitespace and update cache
	token := strings.TrimSpace(string(data))
	b.lastModTime = stat.ModTime()
	if token != "" {
		b.cachedToken = token
	}

	return b.cachedToken
}

var (
	providerMutex        sync.Mutex
	cachedProvider       *oidc.Provider
	oauth2Config         *oauth2.Config
	idTokenVerifier      *oidc.IDTokenVerifier
	ingestProviderMutex  sync.Mutex
	cachedIngestProvider *oidc.Provider
)

// getHTTPClient builds an HTTP client with custom TLS configuration and CA bundle.
func getHTTPClient() (*http.Client, error) {
	config := configuration.GetConfiguration()
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Load system cert pool
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system cert pool: %w", err)
	}

	// Append custom CA file if provided
	if config.OIDCCAFile != "" {
		caCert, err := os.ReadFile(config.OIDCCAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}
		if !rootCAs.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
	}

	tlsConfig.RootCAs = rootCAs
	// Clone default transport and customize
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}, nil
}

// getIngestHTTPClient builds an HTTP client for ingestion OIDC with custom CA bundle.
func getIngestHTTPClient() (*http.Client, error) {
	config := configuration.GetConfiguration()
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Load system cert pool
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system cert pool: %w", err)
	}

	// Append custom CA file if provided
	if config.IngestOIDCCAFile != "" {
		caCert, err := os.ReadFile(config.IngestOIDCCAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read ingestion CA file: %w", err)
		}
		if !rootCAs.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse ingestion CA certificate")
		}
	}

	tlsConfig.RootCAs = rootCAs
	// Clone default transport and customize
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig

	// If JWKS bearer file is provided, wrap the transport with a bearer token RoundTripper
	var rt http.RoundTripper = transport
	if config.IngestOIDCJWKSBearerFile != "" {
		// Extract issuer host from URL
		issuerURL, err := url.Parse(config.IngestOIDCIssuer)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ingestion OIDC issuer URL: %w", err)
		}

		rt = &bearerTokenRoundTripper{
			rt:                transport,
			issuerHost:        issuerURL.Host,
			tokenFilePath:     config.IngestOIDCJWKSBearerFile,
			insecureAllowHTTP: config.IngestOIDCInsecureAllowHTTP,
		}
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: rt,
	}, nil
}

// GetProvider returns the cached OIDC provider or initializes it.
// It's lazy and thread-safe.
func GetProvider(ctx context.Context) (*oidc.Provider, error) {
	providerMutex.Lock()
	defer providerMutex.Unlock()

	// Double-check after acquiring lock
	if cachedProvider != nil {
		return cachedProvider, nil
	}

	config := configuration.GetConfiguration()
	httpClient, err := getHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Create context with custom HTTP client for discovery
	ctx = oidc.ClientContext(ctx, httpClient)

	provider, err := oidc.NewProvider(ctx, config.OIDCIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC provider: %w", err)
	}

	cachedProvider = provider
	return provider, nil
}

// newOAuth2Config builds the oauth2 config. Without a client secret the client is a
// public client (PKCE only): credentials go in the form body (client_id + code_verifier)
// and no Authorization header is sent. Auto-detect would first try HTTP Basic with an
// empty password, which providers such as Entra ID reject for public clients.
func newOAuth2Config(clientID, clientSecret, redirectURL string, endpoint oauth2.Endpoint, scopes []string) *oauth2.Config {
	if clientSecret == "" {
		endpoint.AuthStyle = oauth2.AuthStyleInParams
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     endpoint,
		Scopes:       scopes,
	}
}

// GetOAuth2Config returns the oauth2 configuration.
func GetOAuth2Config(ctx context.Context) (*oauth2.Config, error) {
	providerMutex.Lock()
	if oauth2Config != nil {
		defer providerMutex.Unlock()
		return oauth2Config, nil
	}
	providerMutex.Unlock()

	// Build config outside the lock to avoid holding it while calling GetProvider
	provider, err := GetProvider(ctx)
	if err != nil {
		return nil, err
	}

	config := configuration.GetConfiguration()
	scopes := []string{oidc.ScopeOpenID}

	// Parse additional scopes
	if config.OIDCScopes != "" {
		scopeList := strings.Split(config.OIDCScopes, ",")
		for _, scope := range scopeList {
			scope = strings.TrimSpace(scope)
			if scope != "" && scope != "openid" {
				scopes = append(scopes, scope)
			}
		}
	} else {
		scopes = append(scopes, "profile", "email")
	}

	cfg := newOAuth2Config(config.OIDCClientID, config.OIDCClientSecret, config.OIDCRedirectURL, provider.Endpoint(), scopes)

	// Store in global while holding lock to prevent double initialization
	providerMutex.Lock()
	if oauth2Config == nil {
		oauth2Config = cfg
	}
	providerMutex.Unlock()

	return oauth2Config, nil
}

// GetIDTokenVerifier returns the ID token verifier.
func GetIDTokenVerifier(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	providerMutex.Lock()
	if idTokenVerifier != nil {
		defer providerMutex.Unlock()
		return idTokenVerifier, nil
	}
	providerMutex.Unlock()

	// Build verifier outside the lock to avoid holding it while calling GetProvider
	provider, err := GetProvider(ctx)
	if err != nil {
		return nil, err
	}

	config := configuration.GetConfiguration()
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.OIDCClientID,
		SupportedSigningAlgs: []string{
			"RS256", "RS384", "RS512",
			"ES256", "ES384", "ES512",
			"PS256", "PS384", "PS512",
		},
	})

	// Store in global while holding lock to prevent double initialization
	providerMutex.Lock()
	if idTokenVerifier == nil {
		idTokenVerifier = verifier
	}
	providerMutex.Unlock()

	return idTokenVerifier, nil
}

// getIngestProvider returns the cached ingestion OIDC provider or initializes it.
// It's lazy and thread-safe, separate from the login provider.
func getIngestProvider(ctx context.Context) (*oidc.Provider, error) {
	ingestProviderMutex.Lock()
	defer ingestProviderMutex.Unlock()

	// Double-check after acquiring lock
	if cachedIngestProvider != nil {
		return cachedIngestProvider, nil
	}

	config := configuration.GetConfiguration()
	httpClient, err := getIngestHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create ingestion HTTP client: %w", err)
	}

	// Create context with custom HTTP client for discovery
	ctx = oidc.ClientContext(ctx, httpClient)

	provider, err := oidc.NewProvider(ctx, config.IngestOIDCIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ingestion OIDC provider: %w", err)
	}

	// The JWKS may live on a different host than the issuer; trust that host too
	// (still subject to the HTTPS-only rule) so the bearer token is attached there.
	if brt, ok := httpClient.Transport.(*bearerTokenRoundTripper); ok {
		var disc struct {
			JWKSURI string `json:"jwks_uri"`
		}
		if err := provider.Claims(&disc); err == nil && disc.JWKSURI != "" {
			if u, err := url.Parse(disc.JWKSURI); err == nil && u.Host != "" {
				brt.setJWKSHost(u.Host)
			}
		}
	}

	cachedIngestProvider = provider
	return provider, nil
}
