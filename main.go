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
	"crypto/subtle"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/api"
	"github.com/falcosecurity/falcosidekick-ui/internal/database"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
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

type dbDefaults struct {
	host     string
	port     int
	name     string
	username string
	password string
	ttl      string
}

var databaseDefaults = map[string]dbDefaults{
	"mysql": {
		host:     "localhost",
		port:     3306,
		name:     "falco",
		username: "root",
		password: "",
		ttl:      "0s",
	},
	"postgres": {
		host:     "localhost",
		port:     5432,
		name:     "falco",
		username: "postgres",
		password: "",
		ttl:      "0s",
	},
	"redis": {
		host:     "localhost",
		port:     6379,
		name:     "",
		username: "",
		password: "",
		ttl:      "0s",
	},
}

const AddEvent = "add-event"

func init() {
	// User interface configuration
	addr := utils.GetStringFlagOrEnvParam("a", "FALCOSIDEKICK_UI_ADDR", "0.0.0.0", "Listen Address")
	port := utils.GetIntFlagOrEnvParam("p", "FALCOSIDEKICK_UI_PORT", 2802, "Listen Port")
	version := flag.Bool("v", false, "Print version")
	dev := utils.GetBoolFlagOrEnvParam("x", "FALCOSIDEKICK_UI_DEV", false, "Allow CORS for development")
	loglevel := utils.GetStringFlagOrEnvParam("l", "FALCOSIDEKICK_UI_LOGLEVEL", "info", "Log Level")
	user := utils.GetStringFlagOrEnvParam("u", "FALCOSIDEKICK_UI_USER", "admin:admin", "User in format <login>:<password>")
	disableauth := utils.GetBoolFlagOrEnvParam("d", "FALCOSIDEKICK_UI_DISABLEAUTH", false, "Disable authentication")

	// Database configuration
	dbBackend := utils.GetStringFlagOrEnvParam("b", "FALCOSIDEKICK_UI_DATABASE_BACKEND", "redis", "Database backend (redis, postgres)")

	// Lookup defaults (fallback to redis if unknown)
	defaults, ok := databaseDefaults[strings.ToLower(*dbBackend)]
	if !ok {
		defaults = databaseDefaults["redis"]
	}

	dbHost := utils.GetStringFlagOrEnvParam("dh", "FALCOSIDEKICK_UI_DB_HOST", defaults.host, "Database host")
	dbPort := utils.GetIntFlagOrEnvParam("dp", "FALCOSIDEKICK_UI_DB_PORT", defaults.port, "Database port")
	dbName := utils.GetStringFlagOrEnvParam("dn", "FALCOSIDEKICK_UI_DB_NAME", defaults.name, "Database name")
	dbUsername := utils.GetStringFlagOrEnvParam("dy", "FALCOSIDEKICK_UI_DB_USERNAME", defaults.username, "Database username")
	dbPassword := utils.GetStringFlagOrEnvParam("dw", "FALCOSIDEKICK_UI_DB_PASSWORD", defaults.password, "Database password")
	dbTTL := utils.GetStringFlagOrEnvParam("dt", "FALCOSIDEKICK_UI_DB_TTL", defaults.ttl, "TTL for Redis keys")

	flag.Usage = func() {
		help := `Usage of Falcosidekick-UI:
-v boolean
      Display version
-a string
      Listen Address (default "0.0.0.0", environment "FALCOSIDEKICK_UI_ADDR")
-p int
      Listen Port (default "2802", environment "FALCOSIDEKICK_UI_PORT")
-l string
      Log level: "debug", "info", "warning", "error" (default "info",  environment "FALCOSIDEKICK_UI_LOGLEVEL")
-u string
      User in format <login>:<password> (default "admin:admin", environment "FALCOSIDEKICK_UI_USER")
-d boolean
      Disable authentication (environment "FALCOSIDEKICK_UI_DISABLEAUTH")
-x boolean
      Allow CORS for development (environment "FALCOSIDEKICK_UI_DEV")
-b string
      Database backend: "redis" or "postgres" or "mysql" (default "redis", environment "FALCOSIDEKICK_UI_DATABASE_BACKEND")
-dh string
	  Database Host for either Redis/MySQL/PostgreSQL (default "localhost", environment "FALCOSIDEKICK_UI_DB_HOST")
-dp string
	  Database Port for Redis/MySQL/PostgreSQL (default "6379", environment "FALCOSIDEKICK_UI_DB_PORT")
-dn string
	  Database Name for MySQL/PostgreSQL (default "falco", environment "FALCOSIDEKICK_UI_DB_NAME")
-dy string
      Database Username for Redis/MySQL/PostgreSQL (default "", environment "FALCOSIDEKICK_UI_DB_USERNAME")
-dw string
      Database Password for Redis/MySQL/PostgreSQL (default "", environment "FALCOSIDEKICK_UI_DB_PASSWORD")
-dt string
      Database TTL for Redis keys, the format is X<unit>,
      with unit (s, m, h, d, W, M, y)" (default "0", environment "FALCOSIDEKICK_UI_DB_TTL")
`
		fmt.Println(help)
	}

	// darkmod := flag.Bool("d", false, "Enable dark mode as default")
	flag.Parse()

	if *version {
		v := configuration.GetVersionInfo()
		fmt.Println(v.String())
		os.Exit(0)
	}

	configuration.CreateConfiguration()
	config := configuration.GetConfiguration()
	if ip := net.ParseIP(*addr); ip == nil {
		utils.WriteLog("fatal", "Failed to parse Listen Address")
	}
	if len(strings.Split(*user, ":")) != 2 {
		*user = "admin:admin"
	}
	config.ListenAddress = *addr
	config.ListenPort = *port

	config.DevMode = *dev
	config.LogLevel = *loglevel
	config.Credentials = *user
	config.DisableAuth = *disableauth

	config.DatabaseBackend = *dbBackend
	config.DbHost = *dbHost
	config.DbPort = *dbPort
	config.DbName = *dbName
	config.DbUsername = *dbUsername
	config.DbPassword = *dbPassword
	config.DbTTL = utils.ConvertToSeconds(*dbTTL)

	if utils.GetPriortiyInt(config.LogLevel) < 0 {
		config.LogLevel = "info"
	}

	// Initialize database
	db, err := database.NewEventDatabase()
	if err != nil {
		utils.WriteLog("fatal", "Failed to initialize database: "+err.Error())
		os.Exit(1)
	}

	// Store database instance for API handlers
	models.SetEventDatabase(db)
	models.CreateOutputs()
}

// @title			Falcosidekick UI
// @version		1.0
// @description	Falcosidekick UI
// @contact.name	Falco Authors
// @contact.url	https://github.com/falcosecurity
// @contact.email	cncf-falco-dev@lists.cncf.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @accept			json
// @produce		json
// @schemes		http
// @host			<your-domain>:2802
// @BasePath		/api/v1
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
	if config.DisableAuth {
		utils.WriteLog("warning", "Auhentication disabled")
		e.Use(middleware.CORS())
	}

	utils.WriteLog("info", fmt.Sprintf("Falcosidekick UI is listening on %v:%v", config.ListenAddress, config.ListenPort))
	utils.WriteLog("info", fmt.Sprintf("Log level is %v", config.LogLevel))

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "docs/")
	})
	e.Static("/*", "frontend/dist").Name = "webui-home"
	e.POST("/", api.AddEvent).Name = AddEvent // for compatibility with old Falcosidekicks

	apiRoute := e.Group("/api/v1")
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
			return false
		},
		Validator: func(username, password string, c echo.Context) (bool, error) {
			config := configuration.GetConfiguration()
			if username == "" || password == "" {
				return true, nil
			}
			if subtle.ConstantTimeCompare([]byte(username+":"+password), []byte(config.Credentials)) == 1 {
				return true, nil
			}
			utils.WriteLog("error", "wrong credentials")
			return false, nil
		},
	}))
	apiRoute.POST("/", api.AddEvent).Name = AddEvent
	apiRoute.POST("/auth", api.Authenticate).Name = "authenticate"
	apiRoute.POST("/authenticate", api.Authenticate).Name = "authenticate"
	apiRoute.GET("/config", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/configuration", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/version", api.GetVersionInfo).Name = "get-version"
	apiRoute.GET("/healthz", api.Healthz).Name = "healthz"
	apiRoute.GET("/outputs", api.GetOutputs).Name = "list-outputs"

	eventsRoute := apiRoute.Group("/events")
	eventsRoute.POST("/add", api.AddEvent).Name = AddEvent
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
