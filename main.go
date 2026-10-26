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

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/api"
	"github.com/falcosecurity/falcosidekick-ui/internal/auth"
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/oidc"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	validator "github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/falcosecurity/falcosidekick-ui/docs"
)

type CustomValidator struct {
	validator *validator.Validate
}

const AddEvent = "add-event"

// Route name constants
const (
	routeNameAuthMe       = "auth-me"
	routeNameAuthenticate = "authenticate"
)

func init() {
	addr := utils.GetStringFlagOrEnvParam("a", "FALCOSIDEKICK_UI_ADDR", "0.0.0.0", "Listen Address")
	redisserver := utils.GetStringFlagOrEnvParam("r", "FALCOSIDEKICK_UI_REDIS_URL", "localhost:6379", "Redis server address")
	redisusername := utils.GetStringFlagOrEnvParam("y", "FALCOSIDEKICK_UI_REDIS_USERNAME", "", "Redis server username")
	redispassword := utils.GetStringFlagOrEnvParam("w", "FALCOSIDEKICK_UI_REDIS_PASSWORD", "", "Redis server password")
	port := utils.GetIntFlagOrEnvParam("p", "FALCOSIDEKICK_UI_PORT", 2802, "Listen Port")
	ttl := utils.GetStringFlagOrEnvParam("t", "FALCOSIDEKICK_UI_TTL", "0s", "TTL for keys, the format is X<unit>, with unit (s, m, h, d, W, M, y)")
	version := flag.Bool("v", false, "Print version")
	dev := utils.GetBoolFlagOrEnvParam("x", "FALCOSIDEKICK_UI_DEV", false, "Allow CORS for development")
	loglevel := utils.GetStringFlagOrEnvParam("l", "FALCOSIDEKICK_UI_LOGLEVEL", "info", "Log Level")
	user := utils.GetStringFlagOrEnvParam("u", "FALCOSIDEKICK_UI_USER", "admin:admin", "User in format <login>:<password>")
	disableauth := utils.GetBoolFlagOrEnvParam("d", "FALCOSIDEKICK_UI_DISABLEAUTH", false, "Disable authentication")
	// OIDC flags
	authMode := utils.GetStringFlagOrEnvParam("auth-mode", "FALCOSIDEKICK_UI_AUTH_MODE", "", "Auth mode: \"basic\"|\"oidc\"|\"none\" (default \"basic\")")
	oidcIssuer := utils.GetStringFlagOrEnvParam("oidc-issuer", "FALCOSIDEKICK_UI_OIDC_ISSUER", "", "OIDC issuer URL")
	oidcClientID := utils.GetStringFlagOrEnvParam("oidc-client-id", "FALCOSIDEKICK_UI_OIDC_CLIENT_ID", "", "OIDC client ID")
	oidcClientSecret := utils.GetStringFlagOrEnvParam("oidc-client-secret", "FALCOSIDEKICK_UI_OIDC_CLIENT_SECRET", "", "OIDC client secret")
	oidcClientSecretFile := utils.GetStringFlagOrEnvParam("oidc-client-secret-file", "FALCOSIDEKICK_UI_OIDC_CLIENT_SECRET_FILE", "", "OIDC client secret file")
	oidcRedirectURL := utils.GetStringFlagOrEnvParam("oidc-redirect-url", "FALCOSIDEKICK_UI_OIDC_REDIRECT_URL", "", "OIDC redirect URL")
	oidcScopes := utils.GetStringFlagOrEnvParam("oidc-scopes", "FALCOSIDEKICK_UI_OIDC_SCOPES", "openid,profile,email", "OIDC scopes (comma-separated)")
	oidcUsernameClaim := utils.GetStringFlagOrEnvParam("oidc-username-claim", "FALCOSIDEKICK_UI_OIDC_USERNAME_CLAIM", "preferred_username", "OIDC username claim")
	oidcGroupsClaim := utils.GetStringFlagOrEnvParam("oidc-groups-claim", "FALCOSIDEKICK_UI_OIDC_GROUPS_CLAIM", "groups", "OIDC groups claim")
	oidcAllowedGroups := utils.GetStringFlagOrEnvParam("oidc-allowed-groups", "FALCOSIDEKICK_UI_OIDC_ALLOWED_GROUPS", "", "OIDC allowed groups (comma-separated)")
	oidcCAFile := utils.GetStringFlagOrEnvParam("oidc-ca-file", "FALCOSIDEKICK_UI_OIDC_CA_FILE", "", "OIDC CA file path")
	oidcInsecureAllowHTTP := utils.GetBoolFlagOrEnvParam("oidc-insecure-allow-http", "FALCOSIDEKICK_UI_OIDC_INSECURE_ALLOW_HTTP", false, "Allow HTTP for OIDC (dev only)")
	oidcPostLogoutRedirect := utils.GetStringFlagOrEnvParam("oidc-post-logout-redirect-url", "FALCOSIDEKICK_UI_OIDC_POST_LOGOUT_REDIRECT_URL", "", "OIDC post-logout redirect URL")
	sessionTTL := utils.GetStringFlagOrEnvParam("session-ttl", "FALCOSIDEKICK_UI_SESSION_TTL", "8h", "Session TTL")
	sessionIdleTimeout := utils.GetStringFlagOrEnvParam("session-idle-timeout", "FALCOSIDEKICK_UI_SESSION_IDLE_TIMEOUT", "1h", "Session idle timeout")

	// Ingestion auth flags
	ingestAuth := utils.GetStringFlagOrEnvParam("ingest-auth", "FALCOSIDEKICK_UI_INGEST_AUTH", "none", "Ingestion auth mode: none|oidc")
	ingestOIDCIssuer := utils.GetStringFlagOrEnvParam("ingest-oidc-issuer", "FALCOSIDEKICK_UI_INGEST_OIDC_ISSUER", "", "Ingestion OIDC issuer URL")
	ingestOIDCAudience := utils.GetStringFlagOrEnvParam("ingest-oidc-audience", "FALCOSIDEKICK_UI_INGEST_OIDC_AUDIENCE", "", "Ingestion OIDC audience")
	ingestOIDCAllowedSubjects := utils.GetStringFlagOrEnvParam("ingest-oidc-allowed-subjects", "FALCOSIDEKICK_UI_INGEST_OIDC_ALLOWED_SUBJECTS", "", "Ingestion OIDC allowed subjects (comma-separated)")
	ingestOIDCAllowedClients := utils.GetStringFlagOrEnvParam("ingest-oidc-allowed-clients", "FALCOSIDEKICK_UI_INGEST_OIDC_ALLOWED_CLIENTS", "", "Ingestion OIDC allowed clients (comma-separated)")
	ingestOIDCRequiredScope := utils.GetStringFlagOrEnvParam("ingest-oidc-required-scope", "FALCOSIDEKICK_UI_INGEST_OIDC_REQUIRED_SCOPE", "", "Ingestion OIDC required scope")
	ingestOIDCCAFile := utils.GetStringFlagOrEnvParam("ingest-oidc-ca-file", "FALCOSIDEKICK_UI_INGEST_OIDC_CA_FILE", "", "Ingestion OIDC CA file path")
	ingestOIDCJWKSBearerFile := utils.GetStringFlagOrEnvParam("ingest-oidc-jwks-bearer-file", "FALCOSIDEKICK_UI_INGEST_OIDC_JWKS_BEARER_FILE", "", "Ingestion OIDC JWKS bearer token file")
	ingestOIDCInsecureAllowHTTP := utils.GetBoolFlagOrEnvParam("ingest-oidc-insecure-allow-http", "FALCOSIDEKICK_UI_INGEST_OIDC_INSECURE_ALLOW_HTTP", false, "Allow HTTP for ingestion OIDC (dev only)")

	flag.Usage = func() {
		help := `Usage of Falcosidekick-UI:
	-a string
	      Listen Address (default "0.0.0.0", environment "FALCOSIDEKICK_UI_ADDR")
	-auth-mode string
	      Auth mode: "basic"|"oidc"|"none" (default "basic", environment "FALCOSIDEKICK_UI_AUTH_MODE")
	-d boolean
	      Disable authentication (environment "FALCOSIDEKICK_UI_DISABLEAUTH")
	-ingest-auth string
	      Ingestion auth mode: "none"|"oidc" (default "none", environment "FALCOSIDEKICK_UI_INGEST_AUTH")
	-ingest-oidc-allowed-clients string
	      Ingestion OIDC allowed clients comma-separated (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_ALLOWED_CLIENTS")
	-ingest-oidc-allowed-subjects string
	      Ingestion OIDC allowed subjects comma-separated (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_ALLOWED_SUBJECTS")
	-ingest-oidc-audience string
	      Ingestion OIDC audience (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_AUDIENCE")
	-ingest-oidc-ca-file string
	      Ingestion OIDC CA file path (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_CA_FILE")
	-ingest-oidc-insecure-allow-http boolean
	      Allow HTTP for ingestion OIDC (dev only, environment "FALCOSIDEKICK_UI_INGEST_OIDC_INSECURE_ALLOW_HTTP")
	-ingest-oidc-issuer string
	      Ingestion OIDC issuer URL (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_ISSUER")
	-ingest-oidc-jwks-bearer-file string
	      Ingestion OIDC JWKS bearer token file (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_JWKS_BEARER_FILE")
	-ingest-oidc-required-scope string
	      Ingestion OIDC required scope (default "", environment "FALCOSIDEKICK_UI_INGEST_OIDC_REQUIRED_SCOPE")
	-l string
	      Log level: "debug", "info", "warning", "error" (default "info",  environment "FALCOSIDEKICK_UI_LOGLEVEL")
	-oidc-allowed-groups string
	      OIDC allowed groups comma-separated (default "", environment "FALCOSIDEKICK_UI_OIDC_ALLOWED_GROUPS")
	-oidc-ca-file string
	      OIDC CA file path (default "", environment "FALCOSIDEKICK_UI_OIDC_CA_FILE")
	-oidc-client-id string
	      OIDC client ID (default "", environment "FALCOSIDEKICK_UI_OIDC_CLIENT_ID")
	-oidc-client-secret string
	      OIDC client secret (default "", environment "FALCOSIDEKICK_UI_OIDC_CLIENT_SECRET")
	-oidc-client-secret-file string
	      OIDC client secret file (default "", environment "FALCOSIDEKICK_UI_OIDC_CLIENT_SECRET_FILE")
	-oidc-groups-claim string
	      OIDC groups claim (default "groups", environment "FALCOSIDEKICK_UI_OIDC_GROUPS_CLAIM")
	-oidc-insecure-allow-http boolean
	      Allow HTTP for OIDC (dev only, environment "FALCOSIDEKICK_UI_OIDC_INSECURE_ALLOW_HTTP")
	-oidc-issuer string
	      OIDC issuer URL (default "", environment "FALCOSIDEKICK_UI_OIDC_ISSUER")
	-oidc-post-logout-redirect-url string
	      OIDC post-logout redirect URL (default "", environment "FALCOSIDEKICK_UI_OIDC_POST_LOGOUT_REDIRECT_URL")
	-oidc-redirect-url string
	      OIDC redirect URL (default "", environment "FALCOSIDEKICK_UI_OIDC_REDIRECT_URL")
	-oidc-scopes string
	      OIDC scopes comma-separated (default "openid,profile,email", environment "FALCOSIDEKICK_UI_OIDC_SCOPES")
	-oidc-username-claim string
	      OIDC username claim (default "preferred_username", environment "FALCOSIDEKICK_UI_OIDC_USERNAME_CLAIM")
	-p int
	      Listen Port (default "2802", environment "FALCOSIDEKICK_UI_PORT")
	-r string
	      Redis server address (default "localhost:6379", environment "FALCOSIDEKICK_UI_REDIS_URL")
	-session-idle-timeout string
	      Session idle timeout (default "1h", environment "FALCOSIDEKICK_UI_SESSION_IDLE_TIMEOUT")
	-session-ttl string
	      Session TTL (default "8h", environment "FALCOSIDEKICK_UI_SESSION_TTL")
	-t string
	      TTL for keys, the format is X<unit>,
	      with unit (s, m, h, d, W, M, y)" (default "0", environment "FALCOSIDEKICK_UI_TTL")
	-u string
	      User in format <login>:<password> (default "admin:admin", environment "FALCOSIDEKICK_UI_USER")
	-v boolean
	      Display version
	-w string
	      Redis password (default "", environment "FALCOSIDEKICK_UI_REDIS_PASSWORD")
	-x boolean
	      Allow CORS for development (environment "FALCOSIDEKICK_UI_DEV")
	-y string
	      Redis username (default "", environment "FALCOSIDEKICK_UI_REDIS_USERNAME")
`
		fmt.Println(help)
	}

	flag.Parse()

	if *version {
		v := configuration.GetVersionInfo()
		fmt.Println(v.String())
		os.Exit(0)
	}

	configuration.CreateConfiguration()
	config := configuration.GetConfiguration()

	// Validate listen address
	if ip := net.ParseIP(*addr); ip == nil {
		utils.WriteLog("fatal", "Failed to parse Listen Address")
	}

	// Resolve and validate auth mode
	resolvedAuthMode, err := configuration.ResolveAuthMode(*authMode, *disableauth)
	if err != nil {
		utils.WriteLog("fatal", err.Error())
	}
	*authMode = resolvedAuthMode

	// Validate basic auth credentials
	if len(strings.Split(*user, ":")) != 2 {
		*user = "admin:admin"
	}

	// Set configuration
	config.ListenAddress = *addr
	config.ListenPort = *port
	config.RedisServer = *redisserver
	config.RedisUsername = *redisusername
	config.RedisPassword = *redispassword
	config.DevMode = *dev
	config.TTL = utils.ConvertToSeconds(*ttl)
	config.LogLevel = *loglevel
	config.Credentials = *user
	config.DisableAuth = *disableauth
	config.AuthMode = *authMode

	// Set DisableAuth when auth mode is "none"
	if config.AuthMode == configuration.AuthModeNone {
		config.DisableAuth = true
	}

	// Set OIDC configuration
	config.OIDCIssuer = *oidcIssuer
	config.OIDCClientID = *oidcClientID
	config.OIDCClientSecret = *oidcClientSecret
	config.OIDCClientSecretFile = *oidcClientSecretFile
	config.OIDCRedirectURL = *oidcRedirectURL
	config.OIDCScopes = *oidcScopes
	config.OIDCUsernameClaim = *oidcUsernameClaim
	config.OIDCGroupsClaim = *oidcGroupsClaim
	config.OIDCAllowedGroups = *oidcAllowedGroups
	config.OIDCCAFile = *oidcCAFile
	config.OIDCInsecureAllowHTTP = *oidcInsecureAllowHTTP
	config.OIDCPostLogoutRedirect = *oidcPostLogoutRedirect
	config.SessionTTL = utils.ConvertToSeconds(*sessionTTL)
	config.SessionIdleTimeout = utils.ConvertToSeconds(*sessionIdleTimeout)

	// Set ingestion auth configuration
	config.IngestAuth = strings.ToLower(strings.TrimSpace(*ingestAuth))

	// Validate ingestion auth mode
	if config.IngestAuth != configuration.AuthModeNone && config.IngestAuth != configuration.AuthModeOIDC {
		utils.WriteLog("fatal", fmt.Sprintf("invalid ingest auth mode: %s (must be 'none' or 'oidc')", *ingestAuth))
	}

	// Validate SESSION_TTL and IDLE_TIMEOUT when OIDC is enabled
	if config.AuthMode == configuration.AuthModeOIDC {
		if config.SessionTTL == 0 {
			utils.WriteLog("fatal", "invalid SESSION_TTL: must be a valid duration (e.g., '8h')")
		}
		if config.SessionIdleTimeout == 0 {
			utils.WriteLog("fatal", "invalid IDLE_TIMEOUT: must be a valid duration (e.g., '30m')")
		}
	}

	config.IngestOIDCIssuer = *ingestOIDCIssuer
	config.IngestOIDCAudience = *ingestOIDCAudience
	config.IngestOIDCAllowedSubjects = *ingestOIDCAllowedSubjects
	config.IngestOIDCAllowedClients = *ingestOIDCAllowedClients
	config.IngestOIDCRequiredScope = *ingestOIDCRequiredScope
	config.IngestOIDCCAFile = *ingestOIDCCAFile
	config.IngestOIDCJWKSBearerFile = *ingestOIDCJWKSBearerFile
	config.IngestOIDCInsecureAllowHTTP = *ingestOIDCInsecureAllowHTTP

	// Read client secret from file if specified
	if config.OIDCClientSecretFile != "" {
		secretBytes, err := os.ReadFile(config.OIDCClientSecretFile)
		if err != nil {
			utils.WriteLog("fatal", fmt.Sprintf("failed to read OIDC client secret file: %v", err))
		}
		config.OIDCClientSecret = strings.TrimSpace(string(secretBytes))
	}

	// Validate log level
	if utils.GetPriortiyInt(config.LogLevel) < 0 {
		config.LogLevel = "info"
	}

	// Validate OIDC configuration if in oidc mode
	if config.AuthMode == configuration.AuthModeOIDC {
		if err := oidc.ValidateConfig(); err != nil {
			utils.WriteLog("fatal", fmt.Sprintf("OIDC configuration error: %v", err))
		}
		if config.OIDCClientSecret == "" {
			utils.WriteLog("info", "OIDC running as a public client with PKCE (no client secret configured)")
		}
	}

	// Validate ingestion auth configuration
	if config.IngestAuth == configuration.AuthModeOIDC {
		if err := oidc.ValidateIngestConfig(); err != nil {
			utils.WriteLog("fatal", fmt.Sprintf("ingestion OIDC configuration error: %v", err))
		}
	}

	// Initialize Redis and models
	client := redis.CreateClient()
	redis.CreateIndex(client)
	models.CreateOutputs()
}

// @title          Falcosidekick UI
// @version        1.0
// @description    Falcosidekick UI
// @contact.name   Falco Authors
// @contact.url    https://github.com/falcosecurity
// @contact.email  cncf-falco-dev@lists.cncf.io
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
// @accept         json
// @produce        json
// @schemes        http
// @host           <your-domain>:2802
// @BasePath       /api/v1
func main() {
	e := echo.New()
	v := &CustomValidator{validator: validator.New()}
	config := configuration.GetConfiguration()

	e.Validator = v
	e.HideBanner = true
	e.HidePort = true

	if config.DevMode {
		utils.WriteLog("warning", "DEV mode enabled")
		e.Use(middleware.CORS())
	}
	if config.DisableAuth || config.AuthMode == configuration.AuthModeNone {
		utils.WriteLog("warning", "Authentication disabled")
		e.Use(middleware.CORS())
	}

	utils.WriteLog("info", fmt.Sprintf("Falcosidekick UI is listening on %v:%v", config.ListenAddress, config.ListenPort))
	utils.WriteLog("info", fmt.Sprintf("Log level is %v", config.LogLevel))
	utils.WriteLog("info", fmt.Sprintf("Auth mode is %v", config.AuthMode))

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "docs/")
	})
	e.Static("/*", "frontend/dist").Name = "webui-home"

	apiRoute := e.Group("/api/v1")

	// Wire up middleware and routes based on auth mode
	switch config.AuthMode {
	case configuration.AuthModeOIDC:
		// OIDC mode: register auth routes and middleware
		oidc.RegisterRoutes(apiRoute)

	case configuration.AuthModeNone:
		// No auth mode: skip all auth
		apiRoute.GET("/auth/me", oidc.Me).Name = routeNameAuthMe
		apiRoute.POST("/auth", api.Authenticate).Name = routeNameAuthenticate
		apiRoute.POST("/authenticate", api.Authenticate).Name = routeNameAuthenticate

	default: // basic mode
		// Basic auth middleware (existing behavior)
		apiRoute.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Skipper: func(c echo.Context) bool {
				if configuration.GetConfiguration().DisableAuth {
					return true
				}
				if c.Request().Method == "POST" {
					return true
				}
				if c.Path() == "/api/v1/healthz" {
					return true
				}
				if c.Path() == "/api/v1/auth/me" {
					return true
				}
				return false
			},
			Validator: auth.ValidateCredentials,
		}))

		apiRoute.GET("/auth/me", oidc.Me).Name = routeNameAuthMe
		apiRoute.POST("/auth", api.Authenticate).Name = routeNameAuthenticate
		apiRoute.POST("/authenticate", api.Authenticate).Name = routeNameAuthenticate
	}

	// Build ingest middleware slice (empty when ingest auth is off)
	var ingestMW []echo.MiddlewareFunc
	if config.IngestAuth == configuration.AuthModeOIDC {
		ingestMW = append(ingestMW, oidc.IngestBearerMiddleware())
	}

	// Register ingestion endpoints exactly ONCE with middleware slice
	e.POST("/", api.AddEvent, ingestMW...).Name = AddEvent
	apiRoute.POST("/", api.AddEvent, ingestMW...).Name = AddEvent

	apiRoute.GET("/config", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/configuration", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/version", api.GetVersionInfo).Name = "get-version"
	apiRoute.GET("/healthz", api.Healthz).Name = "healthz"
	apiRoute.GET("/outputs", api.GetOutputs).Name = "list-outputs"

	eventsRoute := apiRoute.Group("/events")
	eventsRoute.POST("/add", api.AddEvent, ingestMW...).Name = AddEvent
	eventsRoute.GET("/count", api.CountEvent).Name = "count-events"
	eventsRoute.GET("/count/:groupby", api.CountByEvent).Name = "count-events-by"
	eventsRoute.GET("/search", api.Search).Name = "search-keys"

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", config.ListenAddress, config.ListenPort)))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		utils.WriteLog("error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
