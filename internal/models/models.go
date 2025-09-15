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

package models

import (
	"encoding/json"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"
)

// EventDatabase interface is forward declared to avoid circular imports
type EventDatabase interface {
	InsertEvent(event *Event) error
	SearchEvents(args *Arguments) (Results, error)
	CountEvents(args *Arguments) (Results, error)
	CountEventsBy(args *Arguments) (Results, error)
	Close() error
}

var eventDB EventDatabase

// SetEventDatabase sets the global event database instance
func SetEventDatabase(db EventDatabase) {
	eventDB = db
}

// GetEventDatabase returns the global event database instance
func GetEventDatabase() EventDatabase {
	return eventDB
}

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
	Hostname string
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
		Hostname: emptyNull(c.QueryParam("hostname")),
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
