// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

// CountEvents counts total events based on the provided arguments
func CountEvents(db *sql.DB, args *models.Arguments) (models.Results, error) {
	whereClause, params := buildWhereClause(args)

	query := "SELECT COUNT(*) FROM falco_events"
	if whereClause != "" {
		query += " " + whereClause
	}

	var totalCount int64
	err := db.QueryRow(query, params...).Scan(&totalCount)
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
		// MySQL doesn't have UNNEST/string_to_array, but we can use FIND_IN_SET or JSON functions.
		// Assuming tags are stored as comma-separated string.
		// Note: This gives approximate grouping (not true normalization).
		// For better results, tags should be normalized into a separate table.

		additionalCondition := "tags IS NOT NULL AND tags != ''"
		if whereClause != "" {
			countSQL = fmt.Sprintf(`
				SELECT SUBSTRING_INDEX(SUBSTRING_INDEX(tags, ',', numbers.n), ',', -1) AS tag_value,
				       COUNT(*) AS count
				FROM falco_events
				JOIN (
					SELECT 1 n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL
					SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL
					SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9 UNION ALL SELECT 10
				) numbers
				  ON CHAR_LENGTH(tags) - CHAR_LENGTH(REPLACE(tags, ',', '')) >= numbers.n - 1
				%s AND %s
				GROUP BY tag_value
				ORDER BY count DESC
			`, whereClause, additionalCondition)
		} else {
			countSQL = fmt.Sprintf(`
				SELECT SUBSTRING_INDEX(SUBSTRING_INDEX(tags, ',', numbers.n), ',', -1) AS tag_value,
				       COUNT(*) AS count
				FROM falco_events
				JOIN (
					SELECT 1 n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL
					SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL
					SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9 UNION ALL SELECT 10
				) numbers
				  ON CHAR_LENGTH(tags) - CHAR_LENGTH(REPLACE(tags, ',', '')) >= numbers.n - 1
				WHERE %s
				GROUP BY tag_value
				ORDER BY count DESC
			`, additionalCondition)
		}
	} else {
		// Regular field groupings
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
