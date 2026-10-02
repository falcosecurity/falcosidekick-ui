# Falcosidekick-ui

[![Falco Ecosystem Repository](https://github.com/falcosecurity/evolution/blob/main/repos/badges/falco-ecosystem-blue.svg)](https://github.com/falcosecurity/evolution/blob/main/REPOSITORIES.md#ecosystem-scope) [![Incubating](https://img.shields.io/badge/status-incubating-orange?style=for-the-badge)](https://github.com/falcosecurity/evolution/blob/main/REPOSITORIES.md#incubating)


![release](https://flat.badgen.net/github/release/falcosecurity/falcosidekick-ui/latest?color=green) ![last commit](https://flat.badgen.net/github/last-commit/falcosecurity/falcosidekick-ui) ![licence](https://flat.badgen.net/badge/license/Apache/blue) ![docker pulls](https://flat.badgen.net/docker/pulls/falcosecurity/falcosidekick-ui?icon=docker)

## Description

A simple WebUI for displaying latest events from [Falco](https://falco.org). It works as output for [Falcosidekick](https://github.com/falcosecurity/falcosidekick).

## Requirements

Events are stored in a `Redis` server with [`Redisearch`](https://github.com/RediSearch/RediSearch) module (> v2).

**For OIDC mode:** Redis >= 6.2 is required (for `GETDEL` command used in flow token management).

## Usage

### Authentication

The UI supports three authentication modes:

#### Basic Authentication (default)
Traditional HTTP Basic Auth. Set username/password via `-u` flag or `FALCOSIDEKICK_UI_USER` environment variable.

#### OIDC/SSO
OpenID Connect authentication with support for any standard OIDC provider (Keycloak, Dex, Authentik, Okta, Entra ID, Google, etc.).

**OIDC Implementation Notes:**
- Session ID tokens are stored server-side in Redis for secure logout hint delivery (only the session ID is stored in cookies)
- JWKS bearer tokens (when configured) are only sent over HTTPS (or HTTP with insecure-allow-http) to the issuer host and the host of the `jwks_uri` advertised by discovery (these differ on Kubernetes, e.g. issuer `kubernetes.default.svc.cluster.local` vs. API server address); bearer credentials are never sent to any other host
- Ingestion endpoints (`/api/v1/` and `/api/v1/events/add`) remain unauthenticated by default; enable `INGEST_OIDC_*` variables to secure them with OIDC-based bearer token validation

**Enable OIDC:**
```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://idp.example.com \
  -oidc-client-id myapp \
  -oidc-client-secret secret \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback
```

**With internal CA:**
```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://idp.example.com \
  -oidc-client-id myapp \
  -oidc-client-secret secret \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback \
  -oidc-ca-file /etc/ssl/certs/ca-bundle.pem
```

**Keycloak Example:**
```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://keycloak.example.com/realms/myrealm \
  -oidc-client-id falcosidekick-ui \
  -oidc-client-secret <confidential-client-secret> \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback \
  -oidc-scopes openid,profile,email \
  -oidc-username-claim preferred_username \
  -oidc-groups-claim groups \
  -oidc-allowed-groups admins,operators
```

**Public client (no secret):**

If neither `-oidc-client-secret` nor `-oidc-client-secret-file` is set, the UI runs as an OIDC public client: the authorization code flow uses PKCE (S256) and the token request sends `client_id` and `code_verifier` in the form body, with no `Authorization` header and no `client_secret`. A line is logged at startup when this mode is active. With a secret set, behaviour is unchanged (confidential client).

Entra ID steps:
1. App registration > Authentication > Add a platform > **Mobile and desktop applications**, and register the redirect URI (e.g. `https://ui.example.com/api/v1/auth/oidc/callback`).
2. Do not create a client secret.
3. If login fails with `AADSTS7000218` (request body must contain `client_assertion` or `client_secret`), set **Allow public client flows** to **Yes** under Authentication > Advanced settings.

```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://login.microsoftonline.com/<tenant-id>/v2.0 \
  -oidc-client-id <application-client-id> \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback
```

**Kubernetes with Keycloak:**
```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://keycloak.example.com/realms/myrealm \
  -oidc-client-id falcosidekick-ui \
  -oidc-client-secret-file /run/secrets/oidc_client_secret \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback
```

Keycloak realm settings:
- Create a "confidential" client with "Standard flow" and "PKCE method: S256" enabled
- Set valid redirect URI to: `https://ui.example.com/api/v1/auth/oidc/callback`
- Add a groups mapper to include user groups in ID token

**Dex Example:**
Dex is a lightweight OpenID Connect provider ideal for development and testing.

```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://dex.example.com \
  -oidc-client-id falcosidekick-ui \
  -oidc-client-secret <dex-client-secret> \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback \
  -oidc-scopes openid,profile,email
```

Dex configuration (example Dex config.yaml):
```yaml
issuer: https://dex.example.com
storage:
  type: memory
oauth2:
  skipApprovalScreen: true
staticClients:
- id: falcosidekick-ui
  redirectURIs:
  - 'https://ui.example.com/api/v1/auth/oidc/callback'
  name: 'Falcosidekick UI'
  secret: <dex-client-secret>
staticPasswords:
- email: admin@example.com
  hash: $2a$10$N7yBV... (bcrypt hash of password)
  username: admin
  userID: <uuid>
```

#### No Authentication
Disable all authentication:
```bash
./falcosidekick-ui -d
```

#### Event Ingestion Authentication (Optional)
By default, event ingestion endpoints (`POST /`, `POST /api/v1/`, `POST /api/v1/events/add`) are **unauthenticated** and publicly accessible. 
Optional OIDC-based bearer token authentication can be enabled to secure event ingestion from Falcosidekick backends.
This is independent of UI login mode and allows fine-grained access control for event ingestion.

**Disabling event ingestion authentication (default):**
```bash
./falcosidekick-ui -ingest-auth none  # or omit the flag
```
Event ingestion endpoints are accessible without authentication.

**Enabling event ingestion authentication:**
```bash
./falcosidekick-ui -ingest-auth oidc \
  -ingest-oidc-issuer https://idp.example.com \
  -ingest-oidc-audience falcosidekick-backend \
  -ingest-oidc-allowed-subjects falcosidekick-service
```

**Environment Variables:**

| Variable | Required | Default | Description |
| :------- | :------: | :------ | :---------- |
| `FALCOSIDEKICK_UI_INGEST_AUTH` | No | `none` | Set to `oidc` to enable bearer token authentication, `none` for unauthenticated access |
| `FALCOSIDEKICK_UI_INGEST_OIDC_ISSUER` | Yes* | - | OIDC provider issuer URL (must be reachable) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_AUDIENCE` | Yes* | - | Expected audience in token (unique identifier for this integration) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_ALLOWED_SUBJECTS` | No | - | Comma-separated allowlist of subject claims (e.g., `user1,user2`) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_ALLOWED_CLIENTS` | No | - | Comma-separated allowlist of client IDs/authorized parties (e.g., `client1,client2`) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_REQUIRED_SCOPE` | No | - | Required scope name in token (space-separated, checked in `scope` or `scp` claim) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_CA_FILE` | No | - | Path to custom CA certificate for OIDC discovery (PEM format) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_JWKS_BEARER_FILE` | No | - | Path to file containing bearer token for JWKS endpoint access (if protected) |
| `FALCOSIDEKICK_UI_INGEST_OIDC_INSECURE_ALLOW_HTTP` | No | `false` | Allow HTTP (insecure) OIDC issuer URLs for development/testing only (default: HTTPS required) |

**Session Configuration (OIDC mode):**

| Variable | Required | Default | Description |
| :------- | :------: | :------ | :---------- |
| `FALCOSIDEKICK_UI_SESSION_TTL` | No | `28800` (8h) | Session absolute lifetime in seconds (must be > 0 in OIDC mode) |
| `FALCOSIDEKICK_UI_SESSION_IDLE_TIMEOUT` | No | `3600` (1h) | Session idle timeout in seconds; sessions expire if idle for this duration (must be > 0 in OIDC mode) |

*Required only if `FALCOSIDEKICK_UI_INGEST_AUTH=oidc`

**At least one authorization check must be configured** (allowed subjects, allowed clients, or required scope) when `FALCOSIDEKICK_UI_INGEST_AUTH=oidc`.

**Keycloak Client-Credentials Example:**
```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://keycloak.example.com/realms/myrealm \
  -oidc-client-id falcosidekick-ui \
  -oidc-client-secret <client-secret> \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback \
  -ingest-auth oidc \
  -ingest-oidc-issuer https://keycloak.example.com/realms/myrealm \
  -ingest-oidc-audience falcosidekick-backend \
  -ingest-oidc-allowed-clients falcosidekick-backend-service
```

Configure Keycloak:
1. Create a "confidential" client for backend ingestion (e.g., `falcosidekick-backend-service`)
2. Enable "Client Credentials" authentication flow
3. Add an "Audience" mapper to include the audience in the token (set to `falcosidekick-backend`)
4. Generate a client secret and use it in your Falcosidekick backend configuration

**Kubernetes Projected Service Account Token Example:**
```bash
./falcosidekick-ui -auth-mode oidc \
  -oidc-issuer https://kubernetes.default.svc \
  -oidc-client-id falcosidekick-ui \
  -oidc-client-secret <client-secret> \
  -oidc-redirect-url https://ui.example.com/api/v1/auth/oidc/callback \
  -ingest-auth oidc \
  -ingest-oidc-issuer https://kubernetes.default.svc \
  -ingest-oidc-audience falcosidekick-ingestion \
  -ingest-oidc-ca-file /var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
  -ingest-oidc-allowed-subjects system:serviceaccount:falco:falcosidekick
```

Configure Kubernetes:
1. Ensure your cluster has `--service-account-issuer` flag set (typically `https://kubernetes.default.svc`)
2. Verify `system:service-account-issuer-discovery` role is bound to allow public OIDC discovery
3. Create a ServiceAccount in the `falco` namespace for Falcosidekick backend
4. Project its token to the Pod (use `projected` volumes in PodSpec)
5. Configure Falcosidekick backend to use the projected token and send it in the `Authorization: Bearer <token>` header

**Integration with Falcosidekick:**
Configure your Falcosidekick backend to authenticate ingestion requests. Refer to Falcosidekick's WebUI output plugin configuration:
- `webui.oauth2.enabled`: Enable OAuth2 authentication for WebUI output
- `webui.oauth2.issuer`: OIDC issuer URL (should match `FALCOSIDEKICK_UI_INGEST_OIDC_ISSUER`)
- `webui.oauth2.audience`: OAuth2 audience (should match `FALCOSIDEKICK_UI_INGEST_OIDC_AUDIENCE`)
- `webui.tokenfile`: Path to file containing bearer token for ingestion requests

### Options
#### Precedence: flag value -> environment variable value -> default value

```shell
Usage of Falcosidekick-UI:
-a string
      Listen Address (default "0.0.0.0", environment "FALCOSIDEKICK_UI_ADDR")
-d boolean
      Disable authentication (environment "FALCOSIDEKICK_UI_DISABLEAUTH")
-l string
      Log level: "debug", "info", "warning", "error" (default "info",  environment "FALCOSIDEKICK_UI_LOGLEVEL")
-p int
      Listen Port (default "2802", environment "FALCOSIDEKICK_UI_PORT")
-r string
      Redis server address (default "localhost:6379", environment "FALCOSIDEKICK_UI_REDIS_URL")
-t string
      TTL for keys, the format is X<unit>,
      with unit (s, m, h, d, W, M, y)" (default "0", environment "FALCOSIDEKICK_UI_TTL")
-u string
      User in format <login>:<password> (default "admin:admin", environment "FALCOSIDEKICK_UI_USER")
-v boolean
      Display version
-w string
      Redis password (default "", environment "FALCOSIDEKICK_UI_REDIS_PASSWORD")
-x boolean
      Allow CORS for development (environment "FALCOSIDEKICK_UI_DEV")
-y string
      Redis username (default "", environment "FALCOSIDEKICK_UI_REDIS_USERNAME")
```

> If not user is set and the authentication is not disabled, the default user is `admin:admin`

### Run with docker

```shell
docker run -d -p 2802:2802 falcosecurity/falcosidekick-ui
```

### Run

```
git clone https://github.com/falcosecurity/falcosidekick-ui.git
cd falcosidekick-ui

go run .
#or
make falcosidekick-ui && ./falcosidekick-ui
```

### Endpoints

| Route   | Method | Query Parameters | Usage            |
| :------ | :----: | :--------------- | :--------------- |
| `/docs` | `GET`  | none             | Get Swagger Docs |
| `/`     | `GET`  | none             | Display WebUI    |

#### UI

The UI is reachable by default at `http://localhost:2802/`.

#### API

> The prefix for access to the API is `/api/v1/`.
> The base URL for the API is `http://localhost:2802/api/v1/`.

| Route                       | Method | Query Parameters                                                         | Usage                                |
| :-------------------------- | :----: | :----------------------------------------------------------------------- | :----------------------------------- |
| `/`                         | `POST` | none                                                                     | Add event                            |
| `/healthz`                  | `GET`  | none                                                                     | Healthcheck                          |
| `/authenticate`, `/auth`    | `POST` | none                                                                     | Authenticate                         |
| `/configuration`, `/config` | `GET`  | none                                                                     | Get Configuration                    |
| `/outputs`                  | `GET`  | none                                                                     | Get list of Outputs of Falcosidekick |
| `/event/count`              | `GET`  | `pretty`, `priority`, `rule`, `filter`, `tags`, `since`, `limit`, `page` | Count all events                     |
| `/event/count/priority`     | `GET`  | `pretty`, `priority`, `rule`, `filter`, `tags`, `since`, `limit`, `page` | Count events by priority             |
| `/event/count/rule`         | `GET`  | `pretty`, `priority`, `rule`, `filter`, `tags`, `since`, `limit`, `page` | Count events by rule                 |
| `/event/count/source`       | `GET`  | `pretty`, `priority`, `rule`, `filter`, `tags`, `since`, `limit`, `page` | Count events by source               |
| `/event/count/tags`         | `GET`  | `pretty`, `priority`, `rule`, `filter`, `tags`, `since`, `limit`, `page` | Count events by tags                 |
| `/event/search`             | `GET`  | `pretty`, `priority`, `rule`, `filter`, `tags`, `since`, `limit`, `page` | Search events                        |

All responses are in JSON format.

Query parameters list:
* `pretty`: return well formated JSON
* `priority`: filter by priority
* `rule`: filter by rule
* `filter`: filter by term
* `source`: filter by source
* `tags`: filter by tags
* `since`: filter by since (in 'second', 'min', 'day', 'week', 'month', 'year')
* `limit`: limit number of results (default: 100)
* `page`: page of results

## Development

### Start local redis server

```shell
docker run -d -p 6379:6379 redislabs/redisearch:2.2.4
```

### Build

Requirements:
* `go` >= 1.25
* `nodejs` >= v14
* `yarn` >= 1.22

```shell
make falcosidekick-ui
```

### Lint

```shell
make lint
```

### Full lint

```shell
make lint-full
```

### Update Docs

Requirement:
* [`swag`](https://github.com/swaggo/swag)

```shell
make docs
```

## Screenshots

![falcosidekick-ui](imgs/webui_01.png)
![falcosidekick-ui](imgs/webui_02.png)
![falcosidekick-ui](imgs/webui_03.png)
![falcosidekick-ui](imgs/webui_04.png)
![falcosidekick-ui](imgs/webui_05.png)

## Authors

* Thomas Labarussias (https://github.com/Issif)
* Frank Jogeleit (https://github.com/fjogeleit)
