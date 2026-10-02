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

package oidc

import (
	"fmt"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

// ValidateConfig validates OIDC configuration at startup.
func ValidateConfig() error {
	config := configuration.GetConfiguration()

	if config.AuthMode == authModeOIDC {
		if config.OIDCIssuer == "" {
			return fmt.Errorf("OIDC issuer is required in oidc mode")
		}
		if !config.OIDCInsecureAllowHTTP && !strings.HasPrefix(config.OIDCIssuer, "https://") {
			return fmt.Errorf("OIDC issuer must be https:// unless OIDC_INSECURE_ALLOW_HTTP=true")
		}
		if config.OIDCClientID == "" {
			return fmt.Errorf("OIDC client ID is required in oidc mode")
		}
		if config.OIDCRedirectURL == "" {
			return fmt.Errorf("OIDC redirect URL is required in oidc mode")
		}
		if !strings.HasSuffix(config.OIDCRedirectURL, "/api/v1/auth/oidc/callback") {
			return fmt.Errorf("OIDC redirect URL must end with /api/v1/auth/oidc/callback")
		}
		if !config.OIDCInsecureAllowHTTP && !strings.HasPrefix(config.OIDCRedirectURL, "https://") {
			return fmt.Errorf("OIDC redirect URL must be https:// unless OIDC_INSECURE_ALLOW_HTTP=true")
		}
		if config.SessionTTL == 0 {
			return fmt.Errorf("SESSION_TTL is required and must be non-zero in oidc mode")
		}
		if config.SessionIdleTimeout == 0 {
			return fmt.Errorf("SESSION_IDLE_TIMEOUT is required and must be non-zero in oidc mode")
		}
		if config.OIDCInsecureAllowHTTP {
			utils.WriteLog("warning", "OIDC insecure mode enabled - for development only")
		}
	}

	return nil
}

// ValidateIngestConfig validates ingestion OIDC configuration at startup.
func ValidateIngestConfig() error {
	config := configuration.GetConfiguration()

	if config.IngestAuth == authModeOIDC {
		if config.IngestOIDCIssuer == "" {
			return fmt.Errorf("ingestion OIDC issuer is required in oidc ingest mode")
		}
		if !config.IngestOIDCInsecureAllowHTTP && !strings.HasPrefix(config.IngestOIDCIssuer, "https://") {
			return fmt.Errorf("ingestion OIDC issuer must be https:// unless FALCOSIDEKICK_UI_INGEST_OIDC_INSECURE_ALLOW_HTTP=true")
		}
		if config.IngestOIDCAudience == "" {
			return fmt.Errorf("ingestion OIDC audience is required in oidc ingest mode")
		}
		// Prevent using login client ID as ingestion audience
		if config.IngestOIDCAudience == config.OIDCClientID {
			return fmt.Errorf("ingestion OIDC audience must not equal login OIDC client ID")
		}
		// Require at least one authorization check
		if config.IngestOIDCAllowedSubjects == "" &&
			config.IngestOIDCAllowedClients == "" &&
			config.IngestOIDCRequiredScope == "" {
			return fmt.Errorf("at least one of INGEST_OIDC_ALLOWED_SUBJECTS, INGEST_OIDC_ALLOWED_CLIENTS, or INGEST_OIDC_REQUIRED_SCOPE must be set for explicit authorization")
		}
		if config.IngestOIDCInsecureAllowHTTP {
			utils.WriteLog("warning", "ingestion OIDC insecure mode enabled - for development only")
		}
	}

	return nil
}
