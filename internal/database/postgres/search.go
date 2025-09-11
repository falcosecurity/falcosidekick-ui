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

package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

// SearchEvents searches for events based on the provided arguments
func SearchEvents(db *sql.DB, args *models.Arguments) (models.Results, error) {
	whereClause, params := buildWhereClause(args)

	// Count query
	countSQL := "SELECT COUNT(*) FROM falco_events" + whereClause
	var totalCount int64
	err := db.QueryRow(countSQL, params...).Scan(&totalCount)
	if err != nil {
		return models.Results{}, fmt.Errorf("failed to count events: %v", err)
	}

	// Search query with pagination
	limitParam := len(params) + 1
	offsetParam := len(params) + 2
	searchSQL := fmt.Sprintf(`
		SELECT json_data FROM falco_events %s
		ORDER BY timestamp_ns DESC
		LIMIT $%d OFFSET $%d`, whereClause, limitParam, offsetParam) // #nosec G201 - whereClause is built from safe programmatic SQL fragments

	params = append(params, args.Limit, args.Page*args.Limit)

	rows, err := db.Query(searchSQL, params...)
	if err != nil {
		return models.Results{}, fmt.Errorf("failed to search events: %v", err)
	}
	defer rows.Close()

	var events models.Events
	for rows.Next() {
		var jsonData string
		if err := rows.Scan(&jsonData); err != nil {
			return models.Results{}, fmt.Errorf("failed to scan row: %v", err)
		}

		var event models.Event
		if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
			return models.Results{}, fmt.Errorf("failed to unmarshal event: %v", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return models.Results{}, fmt.Errorf("row iteration error: %v", err)
	}

	return models.Results{
		Results: events,
		Stats: models.Statistics{
			All:    totalCount,
			Return: int64(len(events)),
		},
	}, nil
}

// buildWhereClause builds the WHERE clause and parameters for the query
func buildWhereClause(args *models.Arguments) (string, []interface{}) {
	var conditions []string
	var params []interface{}
	paramIndex := 1

	// Filter condition (search in output, rule, etc.)
	if args.Filter != "" {
		conditions = append(conditions, fmt.Sprintf("(output ILIKE $%d OR rule ILIKE $%d)", paramIndex, paramIndex))
		params = append(params, "%"+args.Filter+"%")
		paramIndex++
	}

	// Priority condition
	if args.Priority != "" {
		priorities := strings.Split(args.Priority, ",")
		placeholders := make([]string, len(priorities))
		for i, priority := range priorities {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, strings.TrimSpace(priority))
			paramIndex++
		}
		conditions = append(conditions, fmt.Sprintf("priority IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Rule condition
	if args.Rule != "" {
		rules := strings.Split(args.Rule, ",")
		placeholders := make([]string, len(rules))
		for i, rule := range rules {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, strings.TrimSpace(rule))
			paramIndex++
		}
		conditions = append(conditions, fmt.Sprintf("rule IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Source condition
	if args.Source != "" {
		sources := strings.Split(args.Source, ",")
		placeholders := make([]string, len(sources))
		for i, source := range sources {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, strings.TrimSpace(source))
			paramIndex++
		}
		conditions = append(conditions, fmt.Sprintf("source IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Hostname condition
	if args.Hostname != "" {
		hostnames := strings.Split(args.Hostname, ",")
		placeholders := make([]string, len(hostnames))
		for i, hostname := range hostnames {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, strings.TrimSpace(hostname))
			paramIndex++
		}
		conditions = append(conditions, fmt.Sprintf("hostname IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Tags condition
	if args.Tags != "" {
		tags := strings.Split(args.Tags, ",")
		tagConditions := make([]string, len(tags))
		for i, tag := range tags {
			tagConditions[i] = fmt.Sprintf("tags ILIKE $%d", paramIndex)
			params = append(params, "%"+strings.TrimSpace(tag)+"%")
			paramIndex++
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(tagConditions, " OR ")))
	}

	// Since condition (time-based filtering)
	if args.Since != "" {
		if t := int64(utils.ConvertToSeconds(args.Since)); t != 0 {
			sinceTime := time.Now().UTC().Add(-1 * time.Duration(t) * time.Second)
			conditions = append(conditions, fmt.Sprintf("timestamp_utc >= $%d", paramIndex))
			params = append(params, sinceTime)
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, params
}
