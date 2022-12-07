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

	jsonString, _ := json.Marshal(event)

	doc := redisearch.NewDocument(fmt.Sprintf("event:%v", event.UUID), 1.0).
		Set("rule", event.Rule).
		Set("priority", event.Priority).
		Set("output", event.Output).
		Set("source", event.Source).
		Set("timestamp", event.Time.UnixNano()/1e3).
		Set("tags", strings.Join(event.Tags, ",")).
		Set("json", string(jsonString)).
		Set("uuid", event.UUID).
		SetTTL(c.TTL)
	if event.Hostname != "" {
		doc.Set("hostname", event.Hostname)
	}

	err := client.Index([]redisearch.Document{doc}...)
	return err
}
