# Unreleased

* Add OIDC/SSO authentication support (backend-for-frontend pattern)
  * Support for any standard OIDC provider (Keycloak, Dex, Authentik, Okta, Entra ID, Google, etc.)
  * Session-based authentication with Redis
  * Support for custom CA certificates
  * Group-based access control
* Add OIDC-specific endpoints: `/auth/oidc/login`, `/auth/oidc/callback`, `/auth/logout`, `/auth/me`
* Add new configuration flags for OIDC setup
* Frontend: Support for OIDC login flow with SSO button
* Tokens never reach the browser (BFF pattern)
* Cookie-based sessions with HttpOnly, Secure, SameSite flags
* Add optional OIDC bearer token authentication for event ingestion endpoints
  * RFC 6750 Authorization header Bearer token validation
  * Separate ingestion OIDC provider with custom CA certificate support
  * Optional JWKS bearer token for discovery access
  * Subject, client (azp/client_id), and scope allowlist enforcement
  * Proper WWW-Authenticate header error responses
  * JOSE header typ validation (JWT/at+jwt) with defense-in-depth Keycloak ID token rejection
  * Comprehensive test suite with in-process OIDC issuer

# v2.2.0

* Replace CircleCI by Github Actions
* Add new dialog box to show the details of an alert
* Add an export feature
* Allow to specify the TTL for the keys with an unit (second, minute, day, week)
* Allow to set a password for Redis
* Allow to disable the basic auth
* Add a button to clear the search input
* Fix missing filter/search in query parameters
* Fix search for hostnames and tags with dashes
* Fix error with empty result from `count_by`

# v2.1.0

* Add Basic Auth
* Update log format + remove the bootstrap banner
* Allow to set Log Level
* Add `Hostname` field as filter
* Refactor `Info` view
* Fix the results of `count_by`
* Add autorefresh
* Allow to set a `TTL` for keys in Redis
* Allow to set the settings with env vars
* Use Host URL as API baseURL

# v2.0.2

* Fix force redirect to localhost

# v2.0.1

* Fix empty timelines
* Fix loss of query strings when the page is changed
* Change prefix in logs when an event is added

# v2.0.0

* Full rewrite (frontend + backend)
  * The backend has been rewritten with Echo framework and exposes an API (in RO) to count, search, etc the events.
  * For the storage, and allow full text search, a Redis datababse with module [Redisearch](redis.io/docs/stack/search is used)
  * The frontend is created with Vuejs 2 + Vuetify

# v1.1.0

* Implement `DarkMode` for Falcosidekick UI

# v1.0.1

* Fix `Mixed content` error when the UI is reached by https

# v1.0.0

* New design
* Dynamic time range

# v0.2.0

* Render large output fields as message boxes to prevent horizontal scrolling
* Fix Logo position on larger screens
* Increase the default max number of events in memory to 200 (up from 50)
* Add a frontend based time range filter

# v0.1.0

* First release