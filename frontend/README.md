# Falcosidekick UI

Requirements:
* `nodejs` >= v14
* `yarn` >= 1.22

## Project setup

```
yarn install
```

## Development Prerequiements

* Run Redis docker: `docker run -d -p 6379:6379 redislabs/redisearch:2.2.4`
* Run the Falcosidekick-ui backend: `go run . -x` or `./falcosidekick-ui -x`

### Compiles and hot-reloads for development

```
yarn serve
```

### Compiles and minifies for production

```
yarn build
```

### Lints and fixes files

```
yarn lint
```

###

Access: [`http://localhost:8080`](http://localhost:8080)
