package events

import (
	"net/http"

	redis "github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	models "github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

func Count(a *models.Arguments) (models.Results, error) {
	c := redis.GetClient()
	r, err := redis.CountKey(c, a)
	if err != nil {
		utils.WriteLog("error", err.Error())
		return models.Results{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return r, nil
}

func CountBy(a *models.Arguments) (models.Results, error) {
	c := redis.GetClient()
	r, err := redis.CountKeyBy(c, a)
	if err != nil {
		utils.WriteLog("error", err.Error())
		return models.Results{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return r, nil
}
