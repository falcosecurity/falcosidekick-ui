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

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

// PostgresEventDatabase implements the EventDatabase interface for PostgreSQL
type PostgresEventDatabase struct {
	db *sql.DB
}

// NewPostgresEventDatabase creates a new PostgreSQL event database
func NewPostgresEventDatabase() (*PostgresEventDatabase, error) {
	db := CreateClient()
	if err := CreateSchema(db); err != nil {
		return nil, err
	}
	return &PostgresEventDatabase{
		db: db,
	}, nil
}

// InsertEvent stores a new event in PostgreSQL
func (p *PostgresEventDatabase) InsertEvent(event *models.Event) error {
	return InsertEvent(p.db, event)
}

// SearchEvents retrieves events based on search criteria
func (p *PostgresEventDatabase) SearchEvents(args *models.Arguments) (models.Results, error) {
	return SearchEvents(p.db, args)
}

// CountEvents counts events based on criteria
func (p *PostgresEventDatabase) CountEvents(args *models.Arguments) (models.Results, error) {
	return CountEvents(p.db, args)
}

// CountEventsBy counts events grouped by a field
func (p *PostgresEventDatabase) CountEventsBy(args *models.Arguments) (models.Results, error) {
	return CountEventsBy(p.db, args)
}

// Close closes the PostgreSQL connection
func (p *PostgresEventDatabase) Close() error {
	return Close()
}
