# SPDX-License-Identifier: Apache-2.0
#
# Copyright (C) 2023 The Falco Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
# the License. You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.
#

# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL=/bin/bash -o pipefail

.DEFAULT_GOAL:=help

GOPATH  := $(shell go env GOPATH)
GOARCH  := $(shell go env GOARCH)
GOOS    := $(shell go env GOOS)
GOPROXY := $(shell go env GOPROXY)
ifeq ($(GOPROXY),)
GOPROXY := https://proxy.golang.org
endif
export GOPROXY
GO ?= go
DOCKER ?= docker
TEST_FLAGS ?= -v -race

# Directories.
TOOLS_DIR := tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/bin)
GO_INSTALL = tools/go_install.sh

# Binaries.
GOLANGCI_LINT_VER := v1.52.2
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER)

GIT_TAG ?= dirty-tag
GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GIT_HASH ?= $(shell git rev-parse HEAD)
DATE_FMT = +'%Y-%m-%dT%H:%M:%SZ'
SOURCE_DATE_EPOCH ?= $(shell git log -1 --pretty=%ct)
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
GIT_TREESTATE = "clean"
DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(DIFF), 1)
    GIT_TREESTATE = "dirty"
endif

PKG = github.com/falcosecurity/falcosidekick-ui/configuration
LDFLAGS=-X $(PKG).GitVersion=$(GIT_VERSION) -X $(PKG).GitCommit=$(GIT_HASH) -X $(PKG).GitTreeState=$(GIT_TREESTATE) -X $(PKG).BuildDate=$(BUILD_DATE)

## --------------------------------------
## Dev
## --------------------------------------

.PHONY: server
server:
	cd ./frontend && yarn serve

## --------------------------------------
## Docs
## --------------------------------------

.PHONY: docs
docs:
	swag fmt
	swag init

## --------------------------------------
## Build
## --------------------------------------

.PHONY: frontend
frontend:
	cd frontend \
	&& yarn install \
	&& yarn build

.PHONY: falcosidekick-ui
falcosidekick-ui: frontend
	$(GO) mod download
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o falcosidekick-ui .

.PHONY: falcosidekick-ui-linux-amd64
falcosidekick-ui-linux-amd64: frontend
	$(GO) mod download
	GOOS=linux GOARCH=amd64 $(GO) build -gcflags all=-trimpath=/src -asmflags all=-trimpath=/src -a -installsuffix cgo -o falcosidekick-ui .

.PHONY: falcosidekick-ui-backend-only
falcosidekick-ui-backend-only:
	$(GO) mod download
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o falcosidekick-ui .

.PHONY: falcosidekick-ui-linux-amd64-backend-only
falcosidekick-ui-linux-amd64:
	$(GO) mod download
	GOOS=linux GOARCH=amd64 $(GO) build -gcflags all=-trimpath=/src -asmflags all=-trimpath=/src -a -installsuffix cgo -o falcosidekick-ui .

.PHONY: build-image
build-image:
	$(DOCKER) build -t falcosecurity/falcosidekick-ui:latest . --no-cache

## --------------------------------------
## Test
## --------------------------------------

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...

.PHONY: test-coverage
test-coverage:
	$(GO) test ./outputs -count=1 -cover -v ./...

## --------------------------------------
## Linting
## --------------------------------------

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v

lint-full: $(GOLANGCI_LINT) ## Run slower linters to detect possible issues
	$(GOLANGCI_LINT) run -v --fast=false

## --------------------------------------
## Release
## --------------------------------------

.PHONY: goreleaser
goreleaser: ## Release using goreleaser
	LDFLAGS="$(LDFLAGS)" goreleaser release --clean

.PHONY: goreleaser-snapshot
goreleaser-snapshot: ## Release snapshot using goreleaser
	LDFLAGS="$(LDFLAGS)" goreleaser --snapshot --skip-sign --clean

## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(GOLANGCI_LINT): ## Build golangci-lint from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

.PHONY: clean
clean:
	rm -rf hack/tools/bin
	rm -rf dist