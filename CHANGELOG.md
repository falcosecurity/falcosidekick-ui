# v2.0.1

* Fix empty timelines
* Fix loss of query strings when the page is changed
* Change prefix in logs when an event is added

# v2.0.0

* Full rewrite (frontend + backend)
  * The backend has been rewritten with Echo framework and exposes an API (in RO) to count, search, etc the events.
  * For the storage, and allow full text search, a Redis datababse with module [Redisearch](redis.io/docs/stack/search is used
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