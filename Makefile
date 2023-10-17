.PHONY: help
.DEFAULT_GOAL := help
SHELL=bash

UID := $(shell id -u)

dc_path=./docker-compose.yaml
dc_grafana_path=./docker-compose.grafana.yaml
dev_clickhouse_server_container=dev_clickhouse_server
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Starts the environment
	docker-compose -f ${dc_path} build
	docker-compose -f ${dc_path} up -d

up_with_grafana: ## Starts the environment with Grafana
	make up
	mkdir -p dev/grafana
	docker-compose -f ${dc_grafana_path} build
	env UID=${UID} docker-compose -f ${dc_grafana_path} up

down: ## Stops the environment
	docker-compose -f ${dc_path} down
	env UID=${UID} docker-compose -f ${dc_grafana_path} down

down/clean: down
	rm -rf _dev/clickhouse
	mkdir -p _dev/clickhouse

restart: down up # Restart the environment

rebuild: ## Rebuilds the environment from scratch
	@/bin/echo -n "All the volumes will be deleted. You will loose data in DB. Are you sure? [y/N]: " && read answer && \
	[[ $${answer:-N} = y ]] && make destroy

destroy: ## Destroys thel environment
	docker-compose -f ${dc_path} down --volumes --remove-orphans
	docker-compose -f ${dc_path} rm -vsf
	docker-compose -f ${dc_grafana_path} down --volumes --remove-orphans
	docker-compose -f ${dc_grafana_path} rm -vsf

jump_clickhouse_server: ## Jump to the clickhouse_server container
	docker-compose -f ${dc_path} exec ${dev_clickhouse_server_container} bash

.PHONY: lint/install
lint/install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: mockgen/install
mockgen/install:
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: tools/install
tools/install: lint/install mockgen/install

.PHONY: test
test:
	BRUM_SERVER_HOST=localhost \
	BRUM_SERVER_PORT=8087 \
	BRUM_DATABASE_HOST=localhost \
	BRUM_DATABASE_PORT=9000 \
	BRUM_DATABASE_NAME=default \
	BRUM_DATABASE_USERNAME=default \
	BRUM_DATABASE_PASSWORD= \
	BRUM_DATABASE_TABLE_PREFIX=integration_test_ \
	BRUM_PERSISTANCE_DATABASE_STRATEGY=all_in_one_db \
	BRUM_PERSISTANCE_TABLE_STRATEGY=all_in_one_table \
	BRUM_BACKUP_ENABLED=false \
	BRUM_BACKUP_DIRECTORY=/home/basicrum_archive \
	BRUM_BACKUP_EXPIRED_DIRECTORY=/home/basicrum_expired \
	BRUM_BACKUP_UNKNOWN_DIRECTORY=/home/basicrum_unknown \
	BRUM_BACKUP_INTERVAL_SECONDS=5 \
	go test --short ./... 

.PHONY: integration
integration:
	SKIP_E2E=true \
	BRUM_SERVER_HOST=localhost \
	BRUM_SERVER_PORT=8087 \
	BRUM_DATABASE_HOST=localhost \
	BRUM_DATABASE_PORT=9000 \
	BRUM_DATABASE_NAME=default \
	BRUM_DATABASE_USERNAME=default \
	BRUM_DATABASE_PASSWORD= \
	BRUM_DATABASE_TABLE_PREFIX=integration_test_ \
	BRUM_PERSISTANCE_DATABASE_STRATEGY=all_in_one_db \
	BRUM_PERSISTANCE_TABLE_STRATEGY=all_in_one_table \
	BRUM_BACKUP_ENABLED=false \
	BRUM_BACKUP_DIRECTORY=/home/basicrum_archive \
	BRUM_BACKUP_EXPIRED_DIRECTORY=/home/basicrum_expired \
	BRUM_BACKUP_UNKNOWN_DIRECTORY=/home/basicrum_unknown \
	BRUM_BACKUP_INTERVAL_SECONDS=5 \
	go test ./...

.PHONY: e2e
e2e:
	BRUM_SERVER_HOST=localhost \
	BRUM_SERVER_PORT=8087 \
	BRUM_DATABASE_HOST=localhost \
	BRUM_DATABASE_PORT=9000 \
	BRUM_DATABASE_NAME=default \
	BRUM_DATABASE_USERNAME=default \
	BRUM_DATABASE_PASSWORD= \
	BRUM_DATABASE_TABLE_PREFIX=integration_test_ \
	BRUM_PERSISTANCE_DATABASE_STRATEGY=all_in_one_db \
	BRUM_PERSISTANCE_TABLE_STRATEGY=all_in_one_table \
	BRUM_BACKUP_ENABLED=false \
	BRUM_BACKUP_DIRECTORY=/home/basicrum_archive \
	BRUM_BACKUP_EXPIRED_DIRECTORY=/home/basicrum_expired \
	BRUM_BACKUP_UNKNOWN_DIRECTORY=/home/basicrum_unknown \
	BRUM_BACKUP_INTERVAL_SECONDS=5 \
	go test ./...
	
.PHONY: docker-unit-test
docker-unit-test:
	docker-compose -f docker-compose.test.yaml up --exit-code-from unit_test unit_test

.PHONY: docker-clean-test
docker-clean-test: 
	docker-compose -f docker-compose.test.yaml down --remove-orphans
	docker-compose -f docker-compose.test-noprefix.yaml down --remove-orphans
	docker-compose -f docker-compose.test-letsencrypt.yaml down --remove-orphans

.PHONY: _docker-integration-test
_docker-integration-test:
	docker-compose -f docker-compose.test.yaml up --build --exit-code-from integration_test integration_test

.PHONY: docker-integration-test
docker-integration-test: _docker-integration-test docker-clean-test

.PHONY: _docker-integration-test-noprefix
_docker-integration-test-noprefix:
	docker-compose -f docker-compose.test-noprefix.yaml up --build --exit-code-from integration_test2 integration_test2

.PHONY: _docker-integration-test-letsencrypt
_docker-integration-test-letsencrypt:
	docker-compose -f docker-compose.test-letsencrypt.yaml up --build --exit-code-from integration_test integration_test

.PHONY: docker-integration-test-noprefix
docker-integration-test-noprefix: _docker-integration-test-noprefix docker-clean-test

.PHONY: docker-integration-test-letsencrypt
docker-integration-test-letsencrypt: _docker-integration-test-letsencrypt docker-clean-test

.PHONY: docker-hub
docker-hub:
	docker build -t basicrum/front_basicrum_go:$(VERSION) .
	# docker push basicrum/front_basicrum_go:$(VERSION)

.PHONY: debug-docker-integration-test
debug-docker-integration-test:
	docker-compose -f docker-compose.test.yaml up --build integration_server

.PHONY: docker/local/build
docker/local/build:
	docker build -t front_basicrum_go .

.PHONY: gen
gen:
	go generate

.PHONY: cover
cover:
	go test -short -cover -coverprofile cover.out ./...
	go tool cover -func=cover.out
	go tool cover -html=cover.out

.PHONY: cover-integration
cover-integration:
	SKIP_E2E=true go test -count=1 -cover -coverprofile cover.out ./...
	go tool cover -func=cover.out
	go tool cover -html=cover.out
