.PHONY: help
help: ## Show the help
	@grep -hE '^[A-Za-z0-9_ \-]*?:.*##.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

GOLANG_VERSION=1.26.0
GOTENBERG_VERSION=8.27.0
APP_NAME=gotenberg-gs
APP_VERSION=8
APP_AUTHOR=Portailpro
APP_REPOSITORY=https://github.com/PortailPro/gotenberg-gs.git
DOCKER_REGISTRY=leblancsimon
DOCKER_REPOSITORY=gotenberg-gs
GHOSTSCRIPT_VERSION=10.06.0

.PHONY: build
build: ## Build the Gotenberg's Docker image
	sudo docker build --no-cache --pull \
	--build-arg GOLANG_VERSION=$(GOLANG_VERSION) \
	--build-arg GOTENBERG_VERSION=$(GOTENBERG_VERSION) \
	--build-arg APP_NAME=$(APP_NAME) \
	--build-arg APP_VERSION=$(APP_VERSION) \
	--build-arg APP_AUTHOR=$(APP_AUTHOR) \
	--build-arg APP_REPOSITORY=$(APP_REPOSITORY) \
	--build-arg GHOSTSCRIPT_VERSION=$(GHOSTSCRIPT_VERSION) \
	-t $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(APP_VERSION) \
	-f build/Dockerfile .

.PHONY: lint
lint: ## Lint Golang codebase
	golangci-lint run

.PHONY: fmt
fmt: ## Format Golang codebase and "optimize" the dependencies
	go fix ./...
	golangci-lint fmt
	go mod tidy
