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

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	_ "github.com/lib/pq"
)

var db *sql.DB

// CreateClient creates a new PostgreSQL connection
func CreateClient() *sql.DB {
	config := configuration.GetConfiguration()

	// Build connection string
	var connStr string
	if config.DbPassword == "" {
		connStr = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			config.DbHost, config.DbPort, config.DbUsername, config.DbName)
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.DbHost, config.DbPort, config.DbUsername, config.DbPassword, config.DbName)
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}

	// Test connection
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping PostgreSQL: %v", err))
	}

	return db
}

// GetClient returns the existing PostgreSQL connection
func GetClient() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("PostgreSQL client not initialized")
	}
	return db, nil
}

// Close closes the PostgreSQL connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
