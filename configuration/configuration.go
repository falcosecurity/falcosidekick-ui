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

package configuration

type Configuration struct {
	// DisplayMode   string `json:"display-mode"`
	ListenAddress string `json:"listen-address"`
	ListenPort    int    `json:"listen-port"`
	RedisServer   string `json:"redis-server"`
	RedisUsername string `json:"redis-username"`
	RedisPassword string `json:"redis-password"`
	DevMode       bool   `json:"dev-mode"`
	DisableAuth   bool   `json:"disable-auth"`
	LogLevel      string `json:"log-level"`
	TTL           int    `json:"ttl"`
	Credentials   string `json:"credentials"`
	// OIDC fields
	AuthMode               string `json:"auth-mode"`
	OIDCIssuer             string `json:"-"`
	OIDCClientID           string `json:"-"`
	OIDCClientSecret       string `json:"-"`
	OIDCClientSecretFile   string `json:"-"`
	OIDCRedirectURL        string `json:"-"`
	OIDCScopes             string `json:"-"`
	OIDCUsernameClaim      string `json:"-"`
	OIDCGroupsClaim        string `json:"-"`
	OIDCAllowedGroups      string `json:"-"`
	OIDCCAFile             string `json:"-"`
	OIDCInsecureAllowHTTP  bool   `json:"-"`
	OIDCPostLogoutRedirect string `json:"-"`
	SessionTTL             int    `json:"-"`
	SessionIdleTimeout     int    `json:"-"`
	// Ingestion auth fields
	IngestAuth                  string `json:"-"`
	IngestOIDCIssuer            string `json:"-"`
	IngestOIDCAudience          string `json:"-"`
	IngestOIDCAllowedSubjects   string `json:"-"`
	IngestOIDCAllowedClients    string `json:"-"`
	IngestOIDCRequiredScope     string `json:"-"`
	IngestOIDCCAFile            string `json:"-"`
	IngestOIDCJWKSBearerFile    string `json:"-"`
	IngestOIDCInsecureAllowHTTP bool   `json:"-"`
}

var config *Configuration

func CreateConfiguration() *Configuration {
	config = new(Configuration)
	return config
}

func GetConfiguration() *Configuration {
	return config
}
