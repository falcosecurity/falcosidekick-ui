package events

import (
	"fmt"
	"net/http"

	"github.com/falcosecurity/falcosidekick-ui/internal/broadcast"
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

func Add(e *models.Event) error {
	client := redis.GetClient()
	if err := redis.SetKey(client, e); err != nil {
		utils.WriteLog("error", err.Error(), false)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	utils.WriteLog("info", fmt.Sprintf("POST event 'event:%v'", e.Time.UnixNano()/1e3), false)

	go broadcast.GetBroadcast().BroadcastMessage()

	return nil
}
