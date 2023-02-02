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

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: docker-unit-test
docker-unit-test:
	docker-compose -f docker-compose.test.yaml up --exit-code-from unit_test unit_test

.PHONY: docker-clean-test
docker-clean-test: 
	docker-compose -f docker-compose.test.yaml down --volumes --remove-orphans
	docker-compose -f docker-compose.test-noprefix.yaml down --volumes --remove-orphans

.PHONY: _docker-integration-test
_docker-integration-test:
	docker-compose -f docker-compose.test.yaml up --exit-code-from integration_test integration_test

.PHONY: docker-integration-test
docker-integration-test: _docker-integration-test docker-clean-test

.PHONY: _docker-integration-test-noprefix
_docker-integration-test-noprefix:
	docker-compose -f docker-compose.test-noprefix.yaml up --exit-code-from integration_test integration_test

.PHONY: docker-integration-test-noprefix
docker-integration-test-noprefix: _docker-integration-test-noprefix docker-clean-test

.PHONY: docker-hub
docker-hub:
	docker build -t basicrum/front_basicrum_go:$(VERSION) .
	# docker push basicrum/front_basicrum_go:$(VERSION)
