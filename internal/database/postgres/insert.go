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

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

const insertEventSQL = `
INSERT INTO falco_events (
    uuid, rule, priority, output, hostname, source, tags, outputfields, 
    timestamp_ns, timestamp_utc, json_data
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) ON CONFLICT (uuid) DO NOTHING
`

// InsertEvent inserts a Falco event into PostgreSQL
func InsertEvent(db *sql.DB, event *models.Event) error {
	// Convert outputfields to string
	of := make([]string, 0, len(event.OutputFields))
	for _, i := range event.OutputFields {
		of = append(of, fmt.Sprintf("%v", i))
	}

	// Convert to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %v", err)
	}

	// Prepare values
	var hostname sql.NullString
	if event.Hostname != "" {
		hostname = sql.NullString{String: event.Hostname, Valid: true}
	}

	timestampNS := event.Time.UnixNano() / 1e3 // Convert to microseconds like Redis

	// Execute insert
	_, err = db.Exec(insertEventSQL,
		event.UUID,
		event.Rule,
		event.Priority,
		event.Output,
		hostname,
		event.Source,
		strings.Join(event.Tags, ","),
		strings.Join(of, ","),
		timestampNS,
		event.Time,
		string(jsonData),
	)

	if err != nil {
		return fmt.Errorf("failed to insert event: %v", err)
	}

	return nil
}
