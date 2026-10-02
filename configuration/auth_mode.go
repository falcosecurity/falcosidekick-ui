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
	"fmt"
	"strings"
)

// Auth mode constants
const (
	AuthModeBasic = "basic"
	AuthModeOIDC  = "oidc"
	AuthModeNone  = "none"
)

// ResolveAuthMode resolves the final auth mode based on raw input and disableAuth flag.
// Logic:
// - raw = strings.ToLower(strings.TrimSpace(raw))
// - disableAuth && raw == "" → "none"
// - disableAuth && raw == "none" → "none"
// - disableAuth && raw is anything else → error (conflict)
// - !disableAuth && raw == "" → "basic"
// - raw in {basic, oidc, none} → raw; otherwise → error "invalid auth mode"
func ResolveAuthMode(raw string, disableAuth bool) (string, error) {
	raw = strings.ToLower(strings.TrimSpace(raw))

	if disableAuth {
		if raw == "" || raw == AuthModeNone {
			return AuthModeNone, nil
		}
		return "", fmt.Errorf("DISABLEAUTH=true conflicts with AUTH_MODE set to something other than 'none'")
	}

	// disableAuth is false
	if raw == "" {
		return AuthModeBasic, nil
	}

	// Validate that raw is one of the valid modes
	if raw == AuthModeBasic || raw == AuthModeOIDC || raw == AuthModeNone {
		return raw, nil
	}

	return "", fmt.Errorf("invalid auth mode: %s", raw)
}
