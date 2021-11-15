package events

import (
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func Search(a *models.Arguments) (models.Results, error) {
	client := redis.GetClient()
	results, err := redis.SearchKey(client, a)
	if err != nil {
		utils.WriteLog("error", err.Error(), false)
		return models.Results{}, err
	}
	return results, nil
}
