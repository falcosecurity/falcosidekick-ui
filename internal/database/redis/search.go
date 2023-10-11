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

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

func SearchKey(client *redisearch.Client, args *models.Arguments) (models.Results, error) {
	results, count, err := client.Search(redisearch.NewQuery(newQuery(args)).SetSortBy("timestamp", false).Limit(args.Page*args.Limit, args.Limit))
	if err != nil {
		return models.Results{}, err
	}

	r := models.Results{}
	evts := models.Events{}
	for _, i := range results {
		var e models.Event
		if err := json.Unmarshal([]byte(i.Properties["json"].(string)), &e); err != nil {
			return models.Results{}, err
		}
		evts = append(evts, e)
	}
	r.Results = evts
	r.Stats.Return = int64(len(results))
	r.Stats.All = int64(count)
	return r, nil
}
