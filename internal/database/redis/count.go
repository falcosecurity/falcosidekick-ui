package redis

import (
	"errors"
	"strconv"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

func CountKey(client *redisearch.Client, args *models.Arguments) (models.Results, error) {
	query := redisearch.NewQuery(newQuery(args))
	_, count, err := client.Aggregate(redisearch.NewAggregateQuery().SetQuery(query))
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
	if (args.GroupBy != "priority" && args.GroupBy != "rule" && args.GroupBy != "source" && args.GroupBy != "tags") || args.GroupBy == "" {
		return models.Results{}, errors.New("wrong Group By field")
	}
	reducer := redisearch.NewReducer("COUNT", []string{})
	groupBy := redisearch.NewGroupBy().AddFields("@" + args.GroupBy).Reduce(*reducer)
	query := redisearch.NewQuery(newQuery(args))

	results, _, err := client.Aggregate(redisearch.NewAggregateQuery().SetQuery(query).GroupBy(*groupBy))
	if err != nil {
		return models.Results{}, err
	}

	ag := make(models.Aggregation)
	var all int64
	for _, i := range results {
		key := i[1]
		if key == "" {
			continue
		}
		keys := strings.Split(key, ",")
		for _, j := range keys {
			count, _ := strconv.ParseInt(i[3], 10, 64)
			ag[j] += count
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
