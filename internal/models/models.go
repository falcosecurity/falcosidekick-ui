package models

import (
	"encoding/json"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"
)

type Events []Event

type Event struct {
	UUID         string                 `json:"uuid" validate:"uuid"`
	Output       string                 `json:"output" validate:"required,ascii"`
	Priority     string                 `json:"priority" validate:"required,printascii"`
	Rule         string                 `json:"rule" validate:"required,printascii"`
	Time         time.Time              `json:"time" validate:"required"`
	Source       string                 `json:"source" validate:"printascii"`
	OutputFields map[string]interface{} `json:"output_fields"`
	Hostname     string                 `json:"hostname" validate:"printascii"`
	Tags         []string               `json:"tags"`
}

type Payload struct {
	Event   Event    `json:"event" validate:"required"`
	Outputs []string `json:"outputs" validate:"required"`
}

type Aggregation map[string]int64

type Results struct {
	Results interface{} `json:"results,omitempty"`
	Stats   Statistics  `json:"statistics,omitempty"`
}

// ResultsCount is used for swagger json only
type ResultsCount struct {
	Stats Statistics `json:"statistics,omitempty"`
}

// ResultsCountBy is used for swagger json only
type ResultsCountBy struct {
	Results map[string]int64 `json:"results,omitempty"`
	Stats   Statistics       `json:"statistics,omitempty"`
}

// ResultsSearch is used for swagger json only
type ResultsSearch struct {
	Results map[int64]Event `json:"results,omitempty"`
	Stats   Statistics      `json:"statistics,omitempty"`
}

type Statistics struct {
	Return   int64 `json:"returned,omitempty"`
	Distinct int64 `json:"distincts,omitempty"`
	All      int64 `json:"all"`
}

type Arguments struct {
	Priority string
	Rule     string
	Since    string
	Source   string
	Tags     string
	Pretty   string
	Page     int
	Limit    int
	Filter   string
	GroupBy  string
}

func GetArguments(c echo.Context) *Arguments {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 100
	}
	page, _ := strconv.Atoi(c.QueryParam("page"))
	page--
	if page < 0 {
		page = 0
	}
	args := &Arguments{
		Priority: emptyNull(c.QueryParam("priority")),
		Rule:     emptyNull(c.QueryParam("rule")),
		Since:    emptyNull(c.QueryParam("since")),
		Source:   emptyNull(c.QueryParam("source")),
		Tags:     emptyNull(c.QueryParam("tags")),
		Pretty:   c.QueryParam("pretty"),
		GroupBy:  c.Param("groupby"),
		Filter:   emptyNull(c.QueryParam("filter")),
		Page:     page,
		Limit:    limit,
	}
	return args
}

func emptyNull(s string) string {
	if s == "null" {
		return ""
	}
	return s
}

func (e Event) String() string {
	j, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(j)
}
