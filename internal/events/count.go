// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package events

import (
	"net/http"

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

func Count(a *models.Arguments) (models.Results, error) {
	db := models.GetEventDatabase()
	if db == nil {
		utils.WriteLog("error", "database not initialized")
		return models.Results{}, echo.NewHTTPError(http.StatusInternalServerError, "database not initialized")
	}

	r, err := db.CountEvents(a)
	if err != nil {
		utils.WriteLog("error", err.Error())
		return models.Results{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return r, nil
}

func CountBy(a *models.Arguments) (models.Results, error) {
	db := models.GetEventDatabase()
	if db == nil {
		utils.WriteLog("error", "database not initialized")
		return models.Results{}, echo.NewHTTPError(http.StatusInternalServerError, "database not initialized")
	}

	r, err := db.CountEventsBy(a)
	if err != nil {
		utils.WriteLog("error", err.Error())
		return models.Results{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return r, nil
}
