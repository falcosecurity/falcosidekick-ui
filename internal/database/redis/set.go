package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Issif/redisearch-go/redisearch"
	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
)

func SetKey(client *redisearch.Client, event *models.Event) error {
	c := configuration.GetConfiguration()
	timestamp := event.Time.UnixNano() / 1e3

	jsonString, _ := json.Marshal(event)

	doc := redisearch.NewDocument(fmt.Sprintf("event:%v", timestamp), 1.0).
		Set("rule", event.Rule).
		Set("priority", event.Priority).
		Set("output", event.Output).
		Set("source", event.Source).
		Set("timestamp", timestamp).
		Set("tags", strings.Join(event.Tags, ",")).
		Set("json", string(jsonString)).
		SetTTL(c.TTL)

	err := client.Index([]redisearch.Document{doc}...)
	return err
}
