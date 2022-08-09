BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
DB_DATA := "$(ROOT_DIR)/db_data"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

run-pgsql-docker:
	docker run --name calendar-db \
	-p 5431:5432 \
	-e POSTGRES_PASSWORD=p@ssw0rd \
	-e POSTGRES_USER=svend \
	-e POSTGRES_DB=db_calendar \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v $(DB_DATA):/var/lib/postgresql/data \
	-d \
	postgres:14.4

stop-pgsql-docker:
	docker stop calendar-db && docker rm calendar-db

run-goose-migration-up:
	cd $(ROOT_DIR)/migrations && \
		goose postgres "user=svend password=p@ssw0rd dbname=db_calendar sslmode=disable port=5431" up

run-goose-migration-down:
	cd $(ROOT_DIR)/migrations && \
		goose postgres "user=svend password=p@ssw0rd dbname=db_calendar sslmode=disable port=5431" down

run-goose-migration-status:
	cd $(ROOT_DIR)/migrations && \
		goose postgres "user=svend password=p@ssw0rd dbname=db_calendar sslmode=disable port=5431" status

migrate: run-goose-migration-up