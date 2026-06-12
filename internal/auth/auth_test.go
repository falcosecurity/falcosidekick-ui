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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	echo "github.com/labstack/echo/v4"
)

const admin = "admin"

func setupConfig(credentials string) {
	configuration.CreateConfiguration()
	config := configuration.GetConfiguration()
	config.Credentials = credentials
}

func TestValidateCredentials(t *testing.T) {
	setupConfig(admin + ":" + admin)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	tests := []struct {
		name     string
		username string
		password string
		want     bool
	}{
		{
			name:     "valid credentials",
			username: admin,
			password: admin,
			want:     true,
		},
		{
			name:     "wrong password",
			username: admin,
			password: "wrong",
			want:     false,
		},
		{
			name:     "wrong username",
			username: "wrong",
			password: admin,
			want:     false,
		},
		{
			name:     "empty password must be rejected",
			username: "anything",
			password: "",
			want:     false,
		},
		{
			name:     "empty username must be rejected",
			username: "",
			password: "anything",
			want:     false,
		},
		{
			name:     "both empty must be rejected",
			username: "",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateCredentials(tt.username, tt.password, c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ValidateCredentials(%q, %q) = %v, want %v", tt.username, tt.password, got, tt.want)
			}
		})
	}
}
