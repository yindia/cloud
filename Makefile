# Project configuration
CLI_NAME := task-cli
SERVER_NAME := task-server
DASHBOARD_NAME := dashboard
BUILD_DIR := bin
CLI_SRC := ./cli
SERVER_SRC := ./server/root
DASHBOARD_SRC := ./clients/dashboard

# Docker configuration
DOCKER_REPO := ghcr.io/bruin-hiring
VERSION := $(shell git describe --tags --always --dirty)
DOCKER_CLI_NAME := task-cli
DOCKER_SERVER_NAME := task-server
DOCKER_DASHBOARD_NAME := task-dashboard

# ANSI color codes for prettier output
NO_COLOR := \033[0m
OK_COLOR := \033[32;01m
ERROR_COLOR := \033[31;01m
WARN_COLOR := \033[33;01m

# Declare phony targets (targets that don't represent files)
.PHONY: all bootstrap deps check-go check-npm build test docker-build docker-push helm-template helm-lint helm-fmt helm-install helm helm-dep-update

# Default target: run deps, tests, and build
all: deps test build

# Install all dependencies
deps: deps-go deps-npm

# Install Go dependencies
deps-go: check-go
	go mod download
	go fmt ./...
	go generate ./...

# Install npm dependencies
deps-npm: check-npm
	npm install --force

# Check if Go is installed
check-go:
	@which go > /dev/null || (echo "$(ERROR_COLOR)Go is not installed$(NO_COLOR)" && exit 1)

# Check if npm is installed
check-npm:
	@which npm > /dev/null || (echo "$(ERROR_COLOR)npm is not installed$(NO_COLOR)" && exit 1)

# CLI targets
build-cli: deps-go
	@echo "$(OK_COLOR)==> Building the CLI...$(NO_COLOR)"
	@CGO_ENABLED=0 go build -v -ldflags="-s -w" -o "$(BUILD_DIR)/$(CLI_NAME)" "$(CLI_SRC)"

run-cli: build-cli
	@echo "$(OK_COLOR)==> Running the CLI...$(NO_COLOR)"
	@$(BUILD_DIR)/$(CLI_NAME) --help

docker-build-cli:
	@echo "$(OK_COLOR)==> Building Docker image for CLI...$(NO_COLOR)"
	docker build --target cli-final -t $(DOCKER_REPO)/$(DOCKER_CLI_NAME):$(VERSION) .

docker-push-cli: docker-build-cli
	@echo "$(OK_COLOR)==> Pushing Docker image for CLI...$(NO_COLOR)"
	docker push $(DOCKER_REPO)/$(DOCKER_CLI_NAME):$(VERSION)

# Server targets
build-server: deps-go
	@echo "$(OK_COLOR)==> Building the server...$(NO_COLOR)"
	@CGO_ENABLED=0 go build -v -ldflags="-s -w" -o "$(BUILD_DIR)/$(SERVER_NAME)" "$(SERVER_SRC)"

run-server: build-server
	@echo "$(OK_COLOR)==> Running the server...$(NO_COLOR)"
	@$(BUILD_DIR)/$(SERVER_NAME)

docker-build-server:
	@echo "$(OK_COLOR)==> Building Docker image for server...$(NO_COLOR)"
	docker build --target server-final -t $(DOCKER_REPO)/$(DOCKER_SERVER_NAME):$(VERSION) .

docker-push-server: docker-build-server
	@echo "$(OK_COLOR)==> Pushing Docker image for server...$(NO_COLOR)"
	docker push $(DOCKER_REPO)/$(DOCKER_SERVER_NAME):$(VERSION)

# Dashboard targets
build-dashboard: deps-npm
	@echo "$(OK_COLOR)==> Building the dashboard...$(NO_COLOR)"
	npm run build

run-dashboard: deps-npm
	@echo "$(OK_COLOR)==> Running the dashboard...$(NO_COLOR)"
	npm run dev

docker-build-dashboard:
	@echo "$(OK_COLOR)==> Building Docker image for dashboard...$(NO_COLOR)"
	docker build -f Dockerfile.client -t $(DOCKER_REPO)/$(DOCKER_DASHBOARD_NAME):$(VERSION) .

docker-push-dashboard: docker-build-dashboard
	@echo "$(OK_COLOR)==> Pushing Docker image for dashboard...$(NO_COLOR)"
	docker push $(DOCKER_REPO)/$(DOCKER_DASHBOARD_NAME):$(VERSION)

# Test targets
test: deps
	@echo "$(OK_COLOR)==> Running the unit tests$(NO_COLOR)"
	@go test -race -tags unit -cover ./...
	cd $(DASHBOARD_SRC) && npm test

# Combined targets
build: build-cli build-server build-dashboard
docker-build: docker-build-cli docker-build-server docker-build-dashboard
docker-push: docker-push-cli docker-push-server docker-push-dashboard

# Helm targets
helm-template:
	@echo "$(OK_COLOR)==> Generating Helm templates...$(NO_COLOR)"
	helm template charts/task

helm-lint:
	@echo "$(OK_COLOR)==> Linting Helm charts...$(NO_COLOR)"
	helm lint charts/task

helm-fmt:
	@echo "$(OK_COLOR)==> Formatting Helm charts...$(NO_COLOR)"
	helm lint --strict charts/task

helm-docs:
	@echo "$(OK_COLOR)==> Generating Helm charts README.md...$(NO_COLOR)"
	go install github.com/norwoodj/helm-docs/cmd/helm-docs@latest
	helm-docs -c  ./charts/task/ 

helm-install:
	@echo "$(OK_COLOR)==> Installing Helm charts...$(NO_COLOR)"
	helm install my-release charts/task

helm-dep-update:
	@echo "$(OK_COLOR)==> Updating Helm dependencies...$(NO_COLOR)"
	helm dependency update ./charts/task/

# Run all Helm-related tasks
helm: helm-dep-update helm-template helm-lint helm-fmt helm-docs
	@echo "$(OK_COLOR)==> Helm template, lint, and format completed.$(NO_COLOR)"

# Set up development environment
bootstrap:
	curl -fsSL https://pixi.sh/install.sh | bash
	brew install bufbuild/buf/buf
	pixi shell