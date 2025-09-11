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

import (
	"fmt"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/database/mysql"
	"github.com/falcosecurity/falcosidekick-ui/internal/database/postgres"
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

// NewEventDatabase creates a new event database based on configuration
func NewEventDatabase() (EventDatabase, error) {
	config := configuration.GetConfiguration()

	// Normalize backend name
	backend := strings.ToLower(strings.TrimSpace(config.DatabaseBackend))
	if backend == "" {
		backend = "redis" // Default to Redis for backward compatibility
	}

	utils.WriteLog("info", fmt.Sprintf("Initializing %s database backend", backend))

	switch backend {
	case "redis":
		return redis.NewRedisEventDatabase(), nil

	case "postgres", "postgresql":
		return postgres.NewPostgresEventDatabase()

	case "mysql":
		return mysql.NewMySQLEventDatabase()

	default:
		return nil, fmt.Errorf("unsupported database backend: %s. Supported backends: redis, postgres", backend)
	}
}
