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

package database

import "github.com/falcosecurity/falcosidekick-ui/internal/models"

// EventDatabase defines the interface for event storage backends
type EventDatabase interface {
	// InsertEvent stores a new event
	InsertEvent(event *models.Event) error

	// SearchEvents retrieves events based on search criteria
	SearchEvents(args *models.Arguments) (models.Results, error)

	// CountEvents counts events based on criteria
	CountEvents(args *models.Arguments) (models.Results, error)

	// CountEventsBy counts events grouped by a field
	CountEventsBy(args *models.Arguments) (models.Results, error)

	// Close closes the database connection
	Close() error
}
