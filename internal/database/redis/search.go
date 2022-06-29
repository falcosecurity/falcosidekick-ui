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
