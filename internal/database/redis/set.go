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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func SetKey(client *redisearch.Client, event *models.Event) error {
	c := configuration.GetConfiguration()

	jsonString, _ := json.Marshal(event)

	of := make([]string, 0, len(event.OutputFields))

	for _, i := range event.OutputFields {
		of = append(of, fmt.Sprintf("%v", i))
	}

	doc := redisearch.NewDocument(fmt.Sprintf("event:%v", event.UUID), 1.0).
		Set("rule", event.Rule).
		Set("priority", event.Priority).
		Set("output", utils.Escape(event.Output)).
		Set("source", event.Source).
		Set("timestamp", event.Time.UnixNano()/1e3).
		Set("tags", utils.Escape(strings.Join(event.Tags, ","))).
		Set("json", string(jsonString)).
		Set("outputfields", utils.Escape(strings.Join(of, ","))).
		Set("uuid", event.UUID).
		SetTTL(c.TTL)
	if event.Hostname != "" {
		doc.Set("hostname", utils.Escape(event.Hostname))
	}

	err := client.Index([]redisearch.Document{doc}...)
	return err
}
