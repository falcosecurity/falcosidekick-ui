package events

import (
	"fmt"
	"net/http"

	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

func Add(e *models.Event) error {
	client := redis.GetClient()
	if err := redis.SetKey(client, e); err != nil {
		utils.WriteLog("error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	utils.WriteLog("debug", fmt.Sprintf("NEW event 'event:%v'", e.UUID))
	return nil
}
