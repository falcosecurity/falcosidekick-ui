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
	"time"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS falco_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) UNIQUE NOT NULL,
    rule TEXT NOT NULL,
    priority VARCHAR(50) NOT NULL,
    output TEXT NOT NULL,
    hostname VARCHAR(255),
    source VARCHAR(255) NOT NULL,
    tags TEXT,
    outputfields TEXT,
    timestamp_ns BIGINT NOT NULL,
    timestamp_utc TIMESTAMP NOT NULL,
    json_data JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

var createIndexStatements = []string{
	`CREATE INDEX idx_falco_events_timestamp_ns ON falco_events(timestamp_ns)`,
	`CREATE INDEX idx_falco_events_timestamp_utc ON falco_events(timestamp_utc)`,
	`CREATE INDEX idx_falco_events_priority ON falco_events(priority)`,
	`CREATE INDEX idx_falco_events_rule ON falco_events(rule(255))`,
	`CREATE INDEX idx_falco_events_source ON falco_events(source)`,
	`CREATE INDEX idx_falco_events_hostname ON falco_events(hostname)`,
	`CREATE INDEX idx_falco_events_uuid ON falco_events(uuid)`,
	// JSON indexing in MySQL requires generated columns. Example:
	// `ALTER TABLE falco_events ADD COLUMN json_rule VARCHAR(255) GENERATED ALWAYS AS (json_unquote(json_extract(json_data, '$.rule'))) STORED`,
	// `CREATE INDEX idx_falco_events_json_rule ON falco_events(json_rule)`,
}

// CreateSchema creates the table and indexes for Falco events
func CreateSchema(db *sql.DB) error {
	// Create table
	if _, err := db.Exec(createTableSQL); err != nil {
		utils.WriteLog("error", "Failed to create falco_events table: "+err.Error())
		return err
	}

	// Create indexes one by one
	for _, indexSQL := range createIndexStatements {
		if _, err := db.Exec(indexSQL); err != nil {
			// Check if it's just a "duplicate key" error (index already exists)
			if !isDuplicateKeyError(err) {
				utils.WriteLog("error", "Failed to create index: "+err.Error())
				return err
			}
		}
	}

	utils.WriteLog("info", "MySQL schema created successfully")
	return nil
}

// isDuplicateKeyError checks if the error is a MySQL duplicate key error
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return errStr == "Error 1061: Duplicate key name" ||
		errStr == "Error 1061 (42000): Duplicate key name" ||
		strings.Contains(errStr, "Duplicate key name")
}

// CleanupOldEvents removes events older than TTL if configured
func CleanupOldEvents(db *sql.DB) error {
	config := configuration.GetConfiguration()
	if config.DbTTL <= 0 {
		return nil // No TTL configured
	}

	cutoffTime := time.Now().Add(-time.Duration(config.DbTTL) * time.Second)

	deleteSQL := `DELETE FROM falco_events WHERE timestamp_utc < ?`
	result, err := db.Exec(deleteSQL, cutoffTime)
	if err != nil {
		utils.WriteLog("error", "Failed to cleanup old events: "+err.Error())
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		utils.WriteLog("info", fmt.Sprintf("Cleaned up %d old events", rowsAffected))
	}

	return nil
}
