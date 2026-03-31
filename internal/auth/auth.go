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

package auth

import (
	"crypto/subtle"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
)

// ValidateCredentials checks username and password against the configured credentials.
// It rejects empty username or password to prevent authentication bypass.
func ValidateCredentials(username, password string, c echo.Context) (bool, error) {
	config := configuration.GetConfiguration()
	if username == "" || password == "" {
		return false, nil
	}
	if subtle.ConstantTimeCompare([]byte(username+":"+password), []byte(config.Credentials)) == 1 {
		return true, nil
	}
	utils.WriteLog("error", "wrong credentials")
	return false, nil
}
