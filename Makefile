PKG_NAME    := github.com/artalkjs/artalk/v2
BIN_NAME	:= ./bin/artalk


HAS_RICHGO  := $(shell which richgo)
GOTEST      ?= $(if $(HAS_RICHGO), richgo test, go test)
ARGS        ?= server

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
VERSION=git-$(subst /,-,$(BRANCH))-$(shell date +%Y%m%d%H)-$(shell git describe --always --dirty)
IMAGE_TAG=$(VERSION)
IAMGE_REPO=d-x.cmstop.net
PROJECT_NAME=artalk

export CGO_ENABLED := 1

all: install build

install:
	go mod tidy

run: all
	$(BIN_NAME) $(ARGS)

build:
	go build \
    	-ldflags "-s -w" \
        -o $(BIN_NAME) \
    	$(PKG_NAME)

linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" go build -tags=jsoniter -v -o ./artalk

ship: linux
	docker build --platform=linux/amd64 -f Dockerfile.server -t ${IAMGE_REPO}/${PROJECT_NAME}:${IMAGE_TAG} .
	docker push ${IAMGE_REPO}/${PROJECT_NAME}:${IMAGE_TAG}
	docker rmi ${IAMGE_REPO}/${PROJECT_NAME}:${IMAGE_TAG}
	docker image prune -f

push: ship

build-frontend:
	./scripts/build-frontend.sh

build-debug:
	@echo "Building Artalk for debugging..."
	@go build \
		-gcflags "all=-N -l" \
		-o $(BIN_NAME) \
		$(PKG_NAME)

dev: build-debug
	ATK_SITE_DEFAULT="ArtalkDocs" \
	ATK_TRUSTED_DOMAINS="http://localhost:5173 http://localhost:23367" \
	$(BIN_NAME) $(ARGS)

test:
	$(GOTEST) -timeout 20m $(or $(TEST_PATHS), ./...)

test-coverage:
	$(GOTEST) -cover $(or $(TEST_PATHS), ./...)

test-coverage-html:
	$(GOTEST) -v -coverprofile=coverage.out $(or $(TEST_PATHS), ./...)
	go tool cover -html=coverage.out

test-frontend-e2e:
	./scripts/frontend-e2e-test.sh $(if $(REPORT), --show-report)

update-i18n:
	go generate ./internal/i18n

update-conf:
	go generate ./internal/config

update-conf-docs:
	go run ./internal/config/meta/gen --format markdown --locale en -o ./docs/docs/en/guide/env.md
	go run ./internal/config/meta/gen --format markdown --locale zh-CN -o ./docs/docs/zh/guide/env.md

update-docs-features:
	pnpm -F docs-landing update:readme

update-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g server/server.go --output ./docs/swagger --requiredByDefault
	pnpm -r swagger:build-http-client

docker-build:
	./scripts/docker-build.sh

docker-push:
	./scripts/docker-build.sh --push

.PHONY: all install run build build-frontend build-debug dev \
	test test-coverage test-coverage-html test-frontend-e2e \
	update-i18n update-conf update-conf-docs \
	update-docs-features update-swagger \
	docker-build docker-push;
