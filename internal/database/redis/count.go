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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func CountKey(client *redisearch.Client, args *models.Arguments) (models.Results, error) {
	query := redisearch.NewQuery(newQuery(args))
	count, _, err := client.AggregateQuery(redisearch.NewAggregateQuery().SetQuery(query))
	if err != nil {
		return models.Results{}, err
	}

	return models.Results{
		Stats: models.Statistics{
			All: int64(count),
		},
	}, nil
}

func CountKeyBy(client *redisearch.Client, args *models.Arguments) (models.Results, error) {
	if (args.GroupBy != "priority" && args.GroupBy != "rule" && args.GroupBy != "source" && args.GroupBy != "tags" && args.GroupBy != "hostname") || args.GroupBy == "" {
		return models.Results{}, errors.New("wrong Group By field")
	}
	reducer := redisearch.NewReducer("COUNT", []string{})
	groupBy := redisearch.NewGroupBy().AddFields("@" + args.GroupBy).Reduce(*reducer)
	query := redisearch.NewQuery(newQuery(args))

	_, results, err := client.AggregateQuery(redisearch.NewAggregateQuery().SetQuery(query).GroupBy(*groupBy))
	if err != nil {
		return models.Results{}, err
	}

	ag := make(models.Aggregation)
	var all int64
	for _, i := range results {
		if i["__generated_aliascount"] == nil {
			return models.Results{}, fmt.Errorf("error result from redis for GroupBy %v", args.GroupBy)
		}
		key := i[args.GroupBy]
		count, _ := strconv.ParseInt(i["__generated_aliascount"].(string), 10, 64)
		if len(key.(string)) == 0 {
			continue
		}
		for _, j := range strings.Split(key.(string), ",") {
			ag[utils.UnEscape(j)] += count
			all += count
		}
	}

	return models.Results{
		Stats: models.Statistics{
			All:      all,
			Distinct: int64(len(ag)),
		},
		Results: ag,
	}, nil
}
