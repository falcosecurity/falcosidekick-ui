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
	"fmt"
	"net"

	"github.com/falcosecurity/falcosidekick-ui/configuration"

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

var client *redisearch.Client

// CreateClient creates a new redisearch.Client in the redis package scope.
func CreateClient() *redisearch.Client {
	config := configuration.GetConfiguration()
	var dialOpts []redis.DialOption

	if config.RedisUsername != "" {
		dialOpts = append(dialOpts, redis.DialUsername(config.RedisUsername))
	}

	if config.RedisPassword != "" {
		dialOpts = append(dialOpts, redis.DialPassword(config.RedisPassword))
	}

	// Validate the host:port address
	host, port, err := net.SplitHostPort(config.RedisServer)
	if err != nil {
		port = "6379"
	}
	if host == "" {
		host = "localhost"
	}
	serverAddress := net.JoinHostPort(host, port)

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", serverAddress, dialOpts...)
		if err != nil {
			return nil, err
		}
		return c, nil
	}}

	client = redisearch.NewClientFromPool(pool, "search-client-1")
	return client
}

// GetClient returns an existing redisearch.Client or an error if the client
// hasn't yet been created.
func GetClient() (*redisearch.Client, error) {
	if client == nil {
		return nil, fmt.Errorf("could not retrieve redisearch.Client")
	}
	return client, nil
}
