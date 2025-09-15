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
	"fmt"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

// CountEvents counts total events based on the provided arguments
func CountEvents(db *sql.DB, args *models.Arguments) (models.Results, error) {
	whereClause, params := buildWhereClause(args)

	countSQL := "SELECT COUNT(*) FROM falco_events" + whereClause
	var totalCount int64
	err := db.QueryRow(countSQL, params...).Scan(&totalCount)
	if err != nil {
		return models.Results{}, fmt.Errorf("failed to count events: %v", err)
	}

	return models.Results{
		Stats: models.Statistics{
			All: totalCount,
		},
	}, nil
}

// CountEventsBy counts events grouped by a specific field
func CountEventsBy(db *sql.DB, args *models.Arguments) (models.Results, error) {
	if args.GroupBy == "" {
		return models.Results{}, fmt.Errorf("GroupBy field is required")
	}

	// Validate GroupBy field
	validFields := map[string]string{
		"priority": "priority",
		"rule":     "rule",
		"source":   "source",
		"hostname": "hostname",
		"tags":     "tags",
	}

	dbField, valid := validFields[args.GroupBy]
	if !valid {
		return models.Results{}, fmt.Errorf("invalid GroupBy field: %s", args.GroupBy)
	}

	whereClause, params := buildWhereClause(args)

	var countSQL string
	if args.GroupBy == "tags" {
		// Special handling for tags since they're comma-separated
		additionalCondition := "tags != ''"
		if whereClause != "" {
			countSQL = fmt.Sprintf(`
				SELECT UNNEST(string_to_array(tags, ',')) as tag_value, COUNT(*) as count
				FROM falco_events 
				%s AND %s
				GROUP BY tag_value
				ORDER BY count DESC
			`, whereClause, additionalCondition)
		} else {
			countSQL = fmt.Sprintf(`
				SELECT UNNEST(string_to_array(tags, ',')) as tag_value, COUNT(*) as count
				FROM falco_events 
				WHERE %s
				GROUP BY tag_value
				ORDER BY count DESC
			`, additionalCondition)
		}
	} else {
		additionalCondition := fmt.Sprintf("%s IS NOT NULL AND %s != ''", dbField, dbField)
		if whereClause != "" {
			countSQL = fmt.Sprintf(`
				SELECT %s, COUNT(*) as count
				FROM falco_events 
				%s AND %s
				GROUP BY %s
				ORDER BY count DESC
			`, dbField, whereClause, additionalCondition, dbField)
		} else {
			countSQL = fmt.Sprintf(`
				SELECT %s, COUNT(*) as count
				FROM falco_events 
				WHERE %s
				GROUP BY %s
				ORDER BY count DESC
			`, dbField, additionalCondition, dbField)
		}
	}

	rows, err := db.Query(countSQL, params...)
	if err != nil {
		return models.Results{}, fmt.Errorf("failed to count events by %s: %v", args.GroupBy, err)
	}
	defer rows.Close()

	aggregation := make(models.Aggregation)
	var totalCount int64

	for rows.Next() {
		var key string
		var count int64
		if err := rows.Scan(&key, &count); err != nil {
			return models.Results{}, fmt.Errorf("failed to scan row: %v", err)
		}

		// Clean up the key value
		key = strings.TrimSpace(key)
		if key != "" {
			aggregation[key] = count
			totalCount += count
		}
	}

	if err := rows.Err(); err != nil {
		return models.Results{}, fmt.Errorf("row iteration error: %v", err)
	}

	return models.Results{
		Results: aggregation,
		Stats: models.Statistics{
			All:      totalCount,
			Distinct: int64(len(aggregation)),
		},
	}, nil
}
