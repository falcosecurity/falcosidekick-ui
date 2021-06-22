# Falcosidekick-ui

![release](https://flat.badgen.net/github/release/falcosecurity/falcosidekick-ui/latest?color=green) ![last commit](https://flat.badgen.net/github/last-commit/falcosecurity/falcosidekick-ui) ![licence](https://flat.badgen.net/badge/license/Apache/blue) ![docker pulls](https://flat.badgen.net/docker/pulls/falcosecurity/falcosidekick-ui?icon=docker) [![falcosidekick-ui](https://circleci.com/gh/falcosecurity/falcosidekick-ui.svg?style=shield)](https://circleci.com/gh/falcosecurity/falcosidekick-ui)

![falcosidekick-ui](https://github.com/falcosecurity/falcosidekick-ui/raw/master/imgs/webui_01.png)
![falcosidekick-ui](https://github.com/falcosecurity/falcosidekick-ui/raw/master/imgs/webui_02.png)

## Description

A simple WebUI for displaying latest events from [Falco](https://falco.org). It works as output for [Falcosidekick](https://github.com/falcosecurity/falcosidekick).

## Usage

### Options

```shell
  -a string
        Listen Address (default "0.0.0.0")
  -p int
        Listen Port (default 2802)
  -r int
        Number of events to keep in retention (default 200)
  -d bool
        Enable dark mode as default
```

### Run with docker

```shell
docker run -d -p 2802:2802 falcosecurity/falcosidekick-ui 
```

### Run

```
git clone https://github.com/falcosecurity/falcosidekick-ui.git
cd ./falcosidekick-ui

go run main.go
#or
make falcosidekick-ui && ./falcosidekick-ui

```

### Endpoint

The UI is reachable by default at `http://localhost:2802/ui`.

### Dark Mode

By default Falcosidekick UI uses the default light theme, you can change this behavior by using the `-d` Flag of the `falcosidekick-ui` CLI. The default is also based on the settings of your browser, if it uses dark mod, Falcosidekick does too. Regardless of the default settings, the theme can be changed at any time in the UI itself.

![falcosidekick-ui](https://github.com/falcosecurity/falcosidekick-ui/raw/master/imgs/webui_03.png)
![falcosidekick-ui](https://github.com/falcosecurity/falcosidekick-ui/raw/master/imgs/webui_04.png)

## Development

### Build

```bash
make falcosidekick-ui
```

### Lint

```bash
make lint
```

Full lint:

```bash
make lint-full
```

## Authors

* Thomas Labarussias (https://github.com/Issif)
* Frank Jogeleit (https://github.com/fjogeleit)

