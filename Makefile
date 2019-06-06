# Build variables
PACKAGE = $(shell go list -m)
BINARY_NAME = $(shell echo ${PACKAGE} | cut -d '/' -f 3)
BUILD_DIR = build
VERSION ?= $(shell git describe --exact-match --tags 2> /dev/null || git rev-parse --short HEAD)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")
LDFLAGS = -ldflags "-w -X $(PACKAGE)/internal/app.Version=${VERSION} -X $(PACKAGE)/internal/app.BuildDate=${BUILD_DATE}"

# Docker variables
DOCKER_IMAGE ?= gggy/wallet
DOCKER_TAG ?= dev

# Dev variables
GO_TEST_PACKAGES = ./pkg/... ./cmd/... ./internal/...
GO_SOURCE_FILES = ./pkg/ ./cmd/ ./internal/


## Templates
define TPL_SVC
.PHONY: build-$1
build-$1: dep
	@CGO_ENABLED=0 go build -ldflags "-w -X $${PACKAGE}/internal/$1/app.Version=${VERSION} -X $${PACKAGE}/internal/$1/app.BuildDate=${BUILD_DATE}" -o $${BUILD_DIR}/$1 $${CMD_DIR}/$1
endef

define TPL_DOCKER
.PHONY: dockerfiles docker-$1
docker-$1:
	@docker build -f ./build/dockerfiles/$1 -t gitlab.sudo.team/nanopool/microsvc/$1:dev .
endef

define TPL_DEPL
.PHONY: dcup-$1
dcup-$1:
	@docker-compose -f ./deployments/$1/docker-compose.local.yml up

.PHONY: dcdown-$1
dcdown-$1:
	@docker-compose -f ./deployments/$1/docker-compose.local.yml down
endef

$(foreach int,$(INTERNALS), $(eval $(call TPL_SVC,$(int))))
$(foreach int,$(INTERNALS), $(eval $(call TPL_DOCKER,$(int))))
$(foreach int,$(INTERNALS), $(eval $(call TPL_DEPL,$(int))))

.PHONY: githooks
githooks: ## Ln githooks
	@rm -rf .git/hooks
	@ln -s ../githooks .git/hooks

.PHONY: dep
dep: ## Install dependencies
	@go mod download
	@go mod vendor
	@go install github.com/golang/mock/mockgen

.PHONE: mock
mock: ## Generate mocks
	@go generate ./...

.PHONY: install-linter
install-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: install-linter
	@golangci-lint run

.PHONY: test
test: dep ## Run tests
	@go test ${ARGS} -v -count=1 -coverprofile .cover ./internal/...

.PHONY: coverage
coverage:
	@go tool cover -html=.cover -o coverage.html

.PHONY: build
build: dep ## Install dependecies and build a binary
	CGO_ENABLED=0 go build -tags '${TAGS}' ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${PACKAGE}/cmd/wallet


.PHONY: docker
docker: ## Build a Docker image
	docker build --rm -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

.PHONY: dcup
dcup:: ## Local docker-compose up
	docker-compose -p $(BINARY_NAME) -f ./deployments/docker-compose.local.yml up

.PHONY: dcdown
dcdown:: ## Local docker-compose down
	docker-compose -p $(BINARY_NAME) -f ./deployments/docker-compose.local.yml down

.PHONY: envcheck
envcheck:: ## Check environment for all the necessary requirements
	$(call executable_check,Go,go)
	$(call executable_check,Docker,docker)
	$(call executable_check,Docker Compose,docker-compose)

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)