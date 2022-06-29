package redis

import (
	"github.com/falcosecurity/falcosidekick-ui/configuration"

	"github.com/Issif/redisearch-go/redisearch"
)

var client *redisearch.Client

func CreateClient() *redisearch.Client {
	config := configuration.GetConfiguration()
	client = redisearch.NewClient(config.RedisServer, "eventIndex")
	return client
}

func GetClient() *redisearch.Client {
	return client
}
