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
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
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

func init() {
	addr := utils.GetStringFlagOrEnvParam("a", "FALCOSIDEKICK_UI_ADDR", "0.0.0.0", "Listen Address")
	redisserver := utils.GetStringFlagOrEnvParam("r", "FALCOSIDEKICK_UI_REDIS_URL", "localhost:6379", "Redis server address")
	redispassword := utils.GetStringFlagOrEnvParam("w", "FALCOSIDEKICK_UI_REDIS_PASSWORD", "", "Redis server password")
	port := utils.GetIntFlagOrEnvParam("p", "FALCOSIDEKICK_UI_PORT", 2802, "Listen Port")
	ttl := utils.GetStringFlagOrEnvParam("t", "FALCOSIDEKICK_UI_TTL", "0s", "TTL for keys, the format is X<unit>, with unit (s, m, h, d, W, M, y)")
	version := flag.Bool("v", false, "Print version")
	dev := utils.GetBoolFlagOrEnvParam("x", "FALCOSIDEKICK_UI_DEV", false, "Allow CORS for development")
	loglevel := utils.GetStringFlagOrEnvParam("l", "FALCOSIDEKICK_UI_LOGLEVEL", "info", "Log Level")
	user := utils.GetStringFlagOrEnvParam("u", "FALCOSIDEKICK_UI_USER", "admin:admin", "User in format <login>:<password>")
	disableauth := utils.GetBoolFlagOrEnvParam("d", "FALCOSIDEKICK_UI_DISABLEAUTH", false, "Disable authentication")

	flag.Usage = func() {
		help := `Usage of Falcosidekick-UI:  
-a string
      Listen Address (default "0.0.0.0", environment "FALCOSIDEKICK_UI_ADDR")
-d boolean
      Disable authentication (environment "FALCOSIDEKICK_UI_DISABLEAUTH")
-l string   
	  Log level: "debug", "info", "warning", "error" (default "info",  environment "FALCOSIDEKICK_UI_LOGLEVEL")
-p int
      Listen Port (default "2802", environment "FALCOSIDEKICK_UI_PORT")
-r string
      Redis server address (default "localhost:6379", environment "FALCOSIDEKICK_UI_REDIS_URL")
-t string
      TTL for keys, the format is X<unit>,
      with unit (s, m, h, d, W, M, y)" (default "0", environment "FALCOSIDEKICK_UI_TTL")
-u string  
      User in format <login>:<password> (default "admin:admin", environment "FALCOSIDEKICK_UI_USER")
-w string  
      Redis password (default "", environment "FALCOSIDEKICK_REDIS_PASSWORD")
-x boolean
      Allow CORS for development (environment "FALCOSIDEKICK_UI_DEV")
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
	config.RedisServer = *redisserver
	config.RedisPassword = *redispassword
	config.DevMode = *dev
	config.TTL = utils.ConvertToSeconds(*ttl)
	config.LogLevel = *loglevel
	config.Credentials = *user
	config.DisableAuth = *disableauth

	if utils.GetPriortiyInt(config.LogLevel) < 0 {
		config.LogLevel = "info"
	}

	redis.CreateClient()
	redis.CreateIndex(redis.GetClient())
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
	c := configuration.GetConfiguration()

	e.Validator = v
	e.HideBanner = true
	e.HidePort = true

	if c.DevMode {
		utils.WriteLog("warning", "DEV mode enabled")
		e.Use(middleware.CORS())
	}
	if c.DisableAuth {
		utils.WriteLog("warning", "Auhentication disabled")
		e.Use(middleware.CORS())
	}
	utils.WriteLog("info", fmt.Sprintf("Falcosidekick UI is listening on %v:%v", c.ListenAddress, c.ListenPort))
	utils.WriteLog("info", fmt.Sprintf("log level is %v", c.LogLevel))

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/docs", func(c echo.Context) error {
		if err := c.Redirect(http.StatusPermanentRedirect, "docs/"); err != nil {
			return err
		}
		return nil
	})
	e.Static("/*", "frontend/dist").Name = "webui-home"
	e.Static("/dashboard", "frontend/dist").Name = "webui-dashboard"
	e.Static("/events", "frontend/dist").Name = "webui-events"
	e.Static("/info", "frontend/dist").Name = "webui-info"
	e.Static("/login", "frontend/dist").Name = "webui-login"
	e.POST("/", api.AddEvent).Name = "add-event" // for compatibility with old Falcosidekicks

	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// }))

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
	apiRoute.POST("/", api.AddEvent).Name = "add-event"
	apiRoute.POST("/auth", api.Authenticate).Name = "authenticate"
	apiRoute.POST("/authenticate", api.Authenticate).Name = "authenticate"
	apiRoute.GET("/config", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/configuration", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/version", api.GetVersionInfo).Name = "get-version"
	apiRoute.GET("/healthz", api.Healthz).Name = "healthz"
	apiRoute.GET("/outputs", api.GetOutputs).Name = "list-outputs"

	eventsRoute := apiRoute.Group("/events")
	eventsRoute.POST("/add", api.AddEvent).Name = "add-event"
	eventsRoute.GET("/count", api.CountEvent).Name = "count-events"
	eventsRoute.GET("/count/:groupby", api.CountByEvent).Name = "count-events-by"
	eventsRoute.GET("/search", api.Search).Name = "search-keys"

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", c.ListenAddress, c.ListenPort)))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		utils.WriteLog("error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
