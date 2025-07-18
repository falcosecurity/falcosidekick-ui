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
}

var config *Configuration

func CreateConfiguration() *Configuration {
	config = new(Configuration)
	return config
}

func GetConfiguration() *Configuration {
	return config
}
