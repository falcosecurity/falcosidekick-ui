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

package redis

import (
	"testing"
)

const (
	defaultRedisPort = "6379"
	defaultRedisHost = "localhost"
	testPort6380     = "6380"
)

// TestResolveRedisAddr tests the resolveRedisAddr helper function.
func TestResolveRedisAddr(t *testing.T) {
	defaultAddr := defaultRedisHost + ":" + defaultRedisPort
	testRedisHost := "redis-master"
	testHostPort := testRedisHost + ":" + defaultRedisPort
	testLocalhostPort := defaultRedisHost + ":" + testPort6380
	testHostPort1 := "h:1"

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "hostname only", input: testRedisHost, expect: testHostPort},
		{name: "empty host with port", input: ":" + testPort6380, expect: testLocalhostPort},
		{name: "host and port", input: testHostPort1, expect: testHostPort1},
		{name: "empty string", input: "", expect: defaultAddr},
		{name: "hostname with default port", input: defaultAddr, expect: defaultAddr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveRedisAddr(tt.input)
			if result != tt.expect {
				t.Errorf("resolveRedisAddr(%q) = %q, want %q", tt.input, result, tt.expect)
			}
		})
	}
}
