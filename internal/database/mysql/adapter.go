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

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

// MySQLEventDatabase implements the EventDatabase interface for MySQL
type MySQLEventDatabase struct {
	db *sql.DB
}

// NewMySQLEventDatabase creates a new MySQL event database
func NewMySQLEventDatabase() (*MySQLEventDatabase, error) {
	db := CreateClient()
	if err := CreateSchema(db); err != nil {
		return nil, err
	}
	return &MySQLEventDatabase{
		db: db,
	}, nil
}

// InsertEvent stores a new event in MySQL
func (m *MySQLEventDatabase) InsertEvent(event *models.Event) error {
	return InsertEvent(m.db, event)
}

// SearchEvents retrieves events based on search criteria
func (m *MySQLEventDatabase) SearchEvents(args *models.Arguments) (models.Results, error) {
	return SearchEvents(m.db, args)
}

// CountEvents counts events based on criteria
func (m *MySQLEventDatabase) CountEvents(args *models.Arguments) (models.Results, error) {
	return CountEvents(m.db, args)
}

// CountEventsBy counts events grouped by a field
func (m *MySQLEventDatabase) CountEventsBy(args *models.Arguments) (models.Results, error) {
	return CountEventsBy(m.db, args)
}

// Close closes the MySQL connection
func (m *MySQLEventDatabase) Close() error {
	return Close()
}
