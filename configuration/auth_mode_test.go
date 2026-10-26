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

import (
	"testing"
)

// TestResolveAuthMode tests the ResolveAuthMode function with table-driven tests.
func TestResolveAuthMode(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		disableAuth bool
		want        string
		wantErr     bool
	}{
		// disableAuth=true cases
		{name: "disableAuth=true, raw empty", raw: "", disableAuth: true, want: AuthModeNone, wantErr: false},
		{name: "disableAuth=true, raw whitespace", raw: "  ", disableAuth: true, want: AuthModeNone, wantErr: false},
		{name: "disableAuth=true, raw none", raw: AuthModeNone, disableAuth: true, want: AuthModeNone, wantErr: false},
		{name: "disableAuth=true, raw NONE (uppercase)", raw: "NONE", disableAuth: true, want: AuthModeNone, wantErr: false},
		{name: "disableAuth=true, raw basic", raw: AuthModeBasic, disableAuth: true, want: "", wantErr: true},
		{name: "disableAuth=true, raw oidc", raw: AuthModeOIDC, disableAuth: true, want: "", wantErr: true},
		{name: "disableAuth=true, raw OIDC (uppercase)", raw: "OIDC", disableAuth: true, want: "", wantErr: true},

		// disableAuth=false cases
		{name: "disableAuth=false, raw empty", raw: "", disableAuth: false, want: AuthModeBasic, wantErr: false},
		{name: "disableAuth=false, raw whitespace", raw: "  ", disableAuth: false, want: AuthModeBasic, wantErr: false},
		{name: "disableAuth=false, raw basic", raw: AuthModeBasic, disableAuth: false, want: AuthModeBasic, wantErr: false},
		{name: "disableAuth=false, raw BASIC (uppercase)", raw: "BASIC", disableAuth: false, want: AuthModeBasic, wantErr: false},
		{name: "disableAuth=false, raw oidc", raw: AuthModeOIDC, disableAuth: false, want: AuthModeOIDC, wantErr: false},
		{name: "disableAuth=false, raw OIDC (uppercase)", raw: "OIDC", disableAuth: false, want: AuthModeOIDC, wantErr: false},
		{name: "disableAuth=false, raw none", raw: AuthModeNone, disableAuth: false, want: AuthModeNone, wantErr: false},
		{name: "disableAuth=false, raw NONE (uppercase)", raw: "NONE", disableAuth: false, want: AuthModeNone, wantErr: false},
		{name: "disableAuth=false, raw invalid", raw: "invalid", disableAuth: false, want: "", wantErr: true},
		{name: "disableAuth=false, raw empty with spaces", raw: "   ", disableAuth: false, want: AuthModeBasic, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAuthMode(tt.raw, tt.disableAuth)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveAuthMode(%q, %v) error = %v, wantErr %v", tt.raw, tt.disableAuth, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveAuthMode(%q, %v) = %q, want %q", tt.raw, tt.disableAuth, got, tt.want)
			}
		})
	}
}
