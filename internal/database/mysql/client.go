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

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// CreateClient creates a new MySQL connection
func CreateClient() *sql.DB {
	config := configuration.GetConfiguration()

	// Build DSN: username:password@tcp(host:port)/dbname?parseTime=true
	var dsn string
	if config.DbPassword == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?parseTime=true",
			config.DbUsername, config.DbHost, config.DbPort, config.DbName)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}

	// Test connection
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping MySQL: %v", err))
	}

	return db
}

// GetClient returns the existing MySQL connection
func GetClient() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("MySQL client not initialized")
	}
	return db, nil
}

// Close closes the MySQL connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
