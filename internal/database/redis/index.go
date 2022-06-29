package redis

import (
	"github.com/Issif/redisearch-go/redisearch"
	utils "github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func isIndexExit(client *redisearch.Client) bool {
	_, err := client.Info()
	if err != nil {
		utils.WriteLog("info", "Index does not exist")
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
		AddField(redisearch.NewTextField("source")).
		AddField(redisearch.NewTextField("tags")).
		AddField(redisearch.NewNumericField("timestamp")).
		AddField(redisearch.NewTextField("json"))

	// Drop an existing index. If the index does not exist an error is returned
	// client.Drop()

	// Create the index with the given schema
	utils.WriteLog("info", "Create Index")
	err := client.CreateIndex(schema)
	utils.CheckErr(err)
}
