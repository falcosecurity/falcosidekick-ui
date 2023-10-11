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
	"github.com/falcosecurity/falcosidekick-ui/internal/database/redis"
	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func Search(a *models.Arguments) (models.Results, error) {
	client := redis.GetClient()
	results, err := redis.SearchKey(client, a)
	if err != nil {
		utils.WriteLog("error", err.Error())
		return models.Results{}, err
	}
	return results, nil
}
