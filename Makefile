B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

GO           ?= go
GOOPTS       ?=
GOOPTS := $(GOOPTS) -mod=vendor
GORACE:=-race
GOBUILD=CGO_ENABLED=0 installsuffix=cgo $(GO) build -ldflags="-X 'main.revision=$(REV)' -s -w" -trimpath $(GOOPTS)
GOTEST=$(GO) test -v $(GORACE) $(GOOPTS)

##@ Lint
.PHONY: lint
lint: ## Runs golangci linter
	@ echo "-> Running linters ..."
	@ golangci-lint run ./... --config ./build/ci/.golangci.yml

beautify: ## Run gofmt, goimports and go mod tidy.
	@echo "\033[2m→ Beautify code...\033[0m"
	gofmt -s -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	go mod tidy

##@ Test
.PHONY: test
test: ## Runs tests for the project (except e2e tests)
	$(GOTEST) $(PKG_LIST)

##@ Build
.PHONY: build
build: info ## Builds secret executable
	$(GOBUILD) -o ./secret ./cmd/secret/main.go

.PHONY: info
info: ## Display build revision
	- @echo "revision $(REV)"

##@ Other
#------------------------------------------------------------------------------
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
#------------- <https://suva.sh/posts/well-documented-makefiles> --------------

.DEFAULT_GOAL := help