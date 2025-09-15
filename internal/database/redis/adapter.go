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

package redis

import (
	"github.com/Issif/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

// RedisEventDatabase implements the EventDatabase interface for Redis
type RedisEventDatabase struct {
	client *redisearch.Client
}

// NewRedisEventDatabase creates a new Redis event database
func NewRedisEventDatabase() *RedisEventDatabase {
	client := CreateClient()
	CreateIndex(client)
	return &RedisEventDatabase{
		client: client,
	}
}

// InsertEvent stores a new event in Redis
func (r *RedisEventDatabase) InsertEvent(event *models.Event) error {
	return SetKey(r.client, event)
}

// SearchEvents retrieves events based on search criteria
func (r *RedisEventDatabase) SearchEvents(args *models.Arguments) (models.Results, error) {
	return SearchKey(r.client, args)
}

// CountEvents counts events based on criteria
func (r *RedisEventDatabase) CountEvents(args *models.Arguments) (models.Results, error) {
	return CountKey(r.client, args)
}

// CountEventsBy counts events grouped by a field
func (r *RedisEventDatabase) CountEventsBy(args *models.Arguments) (models.Results, error) {
	return CountKeyBy(r.client, args)
}

// Close closes the Redis connection
func (r *RedisEventDatabase) Close() error {
	// Redis client doesn't require explicit closing in this implementation
	return nil
}
