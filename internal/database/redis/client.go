package redis

import (
	"github.com/falcosecurity/falcosidekick-ui/configuration"

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

var client *redisearch.Client

func CreateClient() *redisearch.Client {
	config := configuration.GetConfiguration()
	if config.RedisPassword != "" {
		pool := &redis.Pool{Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.RedisServer, redis.DialPassword(config.RedisPassword))
		}}
		client = redisearch.NewClientFromPool(pool, "search-client-1")
	} else {
		client = redisearch.NewClient(config.RedisServer, "eventIndex")
	}
	return client
}

func GetClient() *redisearch.Client {
	return client
}
