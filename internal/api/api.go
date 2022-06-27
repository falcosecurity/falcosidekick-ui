package api

import (
	"fmt"
	"net/http"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/events"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

func returnResult(c echo.Context, r models.Results, err error) error {
	if err != nil {
		utils.WriteLog("error", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if c.QueryParam("pretty") == "true" {
		return c.JSONPretty(http.StatusOK, r, " ")
	}
	return c.JSON(http.StatusOK, r)
}

// Add Event
// @Summary      Add Event
// @Description  Add Event
// @Accept       json
// @Produce      json
// @param        payload  body      models.Payload                  true  "Payload"
// @Success      200      {string}  string                          ""
// @Failure      500      {string}  http.StatusInternalServerError  "Internal Server Error"
// @Router       /api/v1/ [post]
func AddEvent(c echo.Context) error {
	payload := new(models.Payload)
	if err := c.Bind(payload); err != nil {
		utils.WriteLog("error", err.Error())
		err2 := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		if err2 != nil {
			return err
		}
	}

	models.GetOutputs().Update(payload.Outputs)

	if err := events.Add(&payload.Event); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "{}")
}

// Count Events
// @Summary      Count Events
// @Description  Count Events
// @Param        pretty    query  bool    false  "pretty"
// @Param        priority  query  string  false  "priority"
// @Param        source    query  string  false  "source"
// @Param        filter    query  string  false  "filter"
// @Param        rule      query  string  false  "rule"
// @Param        tags      query  string  false  "tags"
// @Param        since     query  int     false  "since"
// @Produce      json
// @Success      200  {object}  models.ResultsCount    "Count Events Result"
// @Failure      400  {string}  http.StatusBadRequest  "Bad Request"
// @Router       /api/v1/events/count [get]
func CountEvent(c echo.Context) error {
	a := models.GetArguments(c)
	utils.WriteLog("info", fmt.Sprintf("GET count (source='%v', priority='%v', rule='%v', since='%v', filter='%v', tags='%v')", a.Source, a.Priority, a.Rule, a.Since, a.Filter, a.Tags))

	r, err := events.Count(a)
	return returnResult(c, r, err)
}

// Count Events By
// @Summary      Count Events By
// @Description  Count Events By
// @Param        groupby   path   string  true   "group By"
// @Param        pretty    query  bool    false  "pretty"
// @Param        priority  query  string  false  "priority"
// @Param        source    query  string  false  "source"
// @Param        filter    query  string  false  "filter"
// @Param        rule      query  string  false  "rule"
// @Param        tags      query  string  false  "tags"
// @Param        since     query  int     false  "since"
// @Produce      json
// @Success      200  {object}  models.ResultsCountBy  "Count Events By Result"
// @Failure      400  {string}  http.StatusBadRequest  "Bad Request"
// @Router       /api/v1/events/count/:groupby [get]
func CountByEvent(c echo.Context) error {
	a := models.GetArguments(c)
	utils.WriteLog("info", fmt.Sprintf("GET count by %v (source='%v', priority='%v', rule='%v', since='%v', filter='%v', tags='%v')", a.GroupBy, a.Source, a.Priority, a.Rule, a.Since, a.Filter, a.Tags))

	r, err := events.CountBy(a)
	return returnResult(c, r, err)
}

// Search Events
// @Summary      Search Events
// @Description  Search Events
// @Param        pretty    query  bool    false  "pretty"
// @Param        priority  query  string  false  "priority"
// @Param        source    query  string  false  "source"
// @Param        filter    query  string  false  "filter"
// @Param        rule      query  string  false  "rule"
// @Param        tags      query  string  false  "tags"
// @Param        since     query  int     false  "since"
// @Produce      json
// @Success      200  {object}  models.ResultsSearch   "Search Events Result"
// @Failure      400  {string}  http.StatusBadRequest  "Bad Request"
// @Router       /api/v1/events/search [get]
func Search(c echo.Context) error {
	a := models.GetArguments(c)
	utils.WriteLog("info", fmt.Sprintf("GET search (source='%v', priority='%v', rule='%v', since='%v', filter='%v', tags='%v', page='%v', limit='%v')", a.Source, a.Priority, a.Rule, a.Since, a.Filter, a.Tags, a.Page, a.Limit))

	r, err := events.Search(a)
	return returnResult(c, r, err)
}

// Healthcheck
// @Summary      Healthcheck
// @Description  Healthcheck
// @Produce      json
// @Success      200  {string}  json  "{\"ok\"}"
// @Router       /api/v1/healthz [get]
// @Failure      500  {string}  http.StatusInternalServerError  "Internal Server Error"
func Healthz(c echo.Context) error {
	msg := make(map[string]string)
	msg["status"] = "ok"
	return c.JSON(http.StatusOK, msg)
}

// List Outputs
// @Summary      List Outputs
// @Description  Healthcheck
// @Produce      json
// @Success      200  {object}  models.Outputs                  "Outputs"
// @Failure      500  {string}  http.StatusInternalServerError  "Internal Server Error"
// @Router       /api/v1/outputs [get]
func GetOutputs(c echo.Context) error {
	utils.WriteLog("info", "GET outputs")
	return c.JSON(http.StatusOK, models.GetOutputs())
}

// Configuration
// @Summary      Configuration
// @Description  Configuration
// @Produce      json
// @Success      200  {object}  configuration.Configuration     "Configuration"
// @Failure      500  {string}  http.StatusInternalServerError  "Internal Server Error"
// @Router       /api/v1/configuration [get]
func GetConfiguration(c echo.Context) error {
	utils.WriteLog("info", "GET config")
	return c.JSON(http.StatusOK, configuration.GetConfiguration())
}

// Version
// @Summary      Version
// @Description  Version
// @Produce      json
// @Success      200  {object}  configuration.VersionInfo       "Version"
// @Failure      500  {string}  http.StatusInternalServerError  "Internal Server Error"
// @Router       /api/v1/version [get]
func GetVersionInfo(c echo.Context) error {
	utils.WriteLog("info", "GET version")
	return c.JSON(http.StatusOK, configuration.GetVersionInfo())
}
