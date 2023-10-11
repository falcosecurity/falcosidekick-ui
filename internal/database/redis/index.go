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
	utils "github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func isIndexExit(client *redisearch.Client) bool {
	_, err := client.Info()
	if err != nil {
		if err.Error() == "NOAUTH Authentication required." {
			utils.WriteLog("fatal", "Authentication required")
		}
		utils.WriteLog("warning", "Index does not exist")
		return false
	}
	return true
}

func CreateIndex(client *redisearch.Client) {
	if isIndexExit(client) {
		return
	}

	// Create a schema
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("output")).
		AddField(redisearch.NewTextField("rule")).
		AddField(redisearch.NewTextField("priority")).
		AddField(redisearch.NewTextField("hostname")).
		AddField(redisearch.NewTextField("source")).
		AddField(redisearch.NewTextField("tags")).
		AddField(redisearch.NewNumericField("timestamp")).
		AddField(redisearch.NewTextField("json"))

	// Drop an existing index. If the index does not exist an error is returned
	// client.Drop()

	// Create the index with the given schema
	utils.WriteLog("warning", "Create Index")
	err := client.CreateIndex(schema)
	utils.CheckErr(err)
}
