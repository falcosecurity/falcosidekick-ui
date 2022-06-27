package redis

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
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
	if (args.GroupBy != "priority" && args.GroupBy != "rule" && args.GroupBy != "source" && args.GroupBy != "tags") || args.GroupBy == "" {
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
		key := i[args.GroupBy]
		count, _ := strconv.ParseInt(i["__generated_aliascount"].(string), 10, 64)
		if len(key.(string)) == 0 {
			continue
		}
		for _, j := range strings.Split(key.(string), ",") {
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
