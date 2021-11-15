package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/api"
	"github.com/falcosecurity/falcosidekick-ui/internal/broadcast"
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	validator "github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/falcosecurity/falcosidekick-ui/docs"
)

// const (
// 	light string = "light"
// 	dark  string = "dark"
// )

// go:embed dist
// var frontend embed.FS

type CustomValidator struct {
	validator *validator.Validate
}

func init() {
	addr := flag.String("a", "0.0.0.0", "Listen Address")
	port := flag.Int("p", 2802, "Listen Port")
	// darkmod := flag.Bool("d", false, "Enable dark mode as default")
	version := flag.Bool("v", false, "Print version")
	dev := flag.Bool("x", false, "Allow CORS for development")
	redisserver := flag.String("r", "localhost:6379", "Redis server address")
	flag.Parse()

	if *version {
		v := configuration.GetVersionInfo()
		fmt.Println(v.String())
		os.Exit(0)
	}

	configuration.CreateConfiguration()
	config := configuration.GetConfiguration()
	if ip := net.ParseIP(*addr); ip == nil {
		utils.WriteLog("error", "Failed to parse Listen Address", true)
	}
	config.ListenAddress = *addr
	config.ListenPort = *port
	// config.DisplayMode = light
	// if *darkmod {
	// 	config.DisplayMode = dark
	// }
	config.RedisServer = *redisserver
	config.DevMode = *dev

	redis.CreateClient()
	redis.CreateIndex(redis.GetClient())
	models.CreateOutputs()
	broadcast.CreateBroadcast()
}

// @title        Falcosidekick UI
// @version      1.0
// @description  Falcosidekick UI

// @contact.name   Falco Authors
// @contact.url    https://github.com/falcosecurity
// @contact.email  cncf-falco-dev@lists.cncf.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @accept   json
// @produce  json

// @schemes  http

// @host      localhost:2802
// @BasePath  /api/v1
func main() {
	e := echo.New()
	v := new(CustomValidator)
	c := configuration.GetConfiguration()

	e.Validator = v

	if c.DevMode {
		utils.WriteLog("info", "DEV mode enabled", false)
		e.Use(middleware.CORS())
	}

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/docs", func(c echo.Context) error { //2
		if err := c.Redirect(http.StatusPermanentRedirect, "docs/"); err != nil {
			return err
		}
		return nil
	})
	e.POST("/", api.AddEvent).Name = "add-event"

	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// }))

	apiRoute := e.Group("/api/v1")
	apiRoute.GET("/configuration", api.GetConfiguration).Name = "get-configuration"
	apiRoute.GET("/version", api.GetVersionInfo).Name = "get-version"
	apiRoute.GET("/healthz", api.Healthz).Name = "healthz"
	apiRoute.GET("/outputs", api.GetOutputs).Name = "list-outputs"
	apiRoute.GET("/config", api.GetConfiguration).Name = "display-cnfig"
	apiRoute.GET("/ws", broadcast.WebSocketBroadcast).Name = "websocket"

	eventsRoute := apiRoute.Group("/events")
	eventsRoute.POST("/add", api.AddEvent).Name = "add-event"
	eventsRoute.GET("/count", api.CountEvent).Name = "count-events"
	eventsRoute.GET("/count/:groupby", api.CountByEvent).Name = "count-events-by"
	eventsRoute.GET("/search", api.Search).Name = "search-keys"

	e.Static("/*", "frontend/dist").Name = "webui-home"
	e.Static("/dashboard", "frontend/dist").Name = "webui-dashboard"
	e.Static("/events", "frontend/dist").Name = "webui-events"
	e.Static("/info", "frontend/dist").Name = "webui-info"

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", c.ListenAddress, c.ListenPort)))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
