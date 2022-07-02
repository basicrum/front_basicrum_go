.PHONY: help
.DEFAULT_GOAL := help
SHELL=bash

dc_path=./docker-compose.yaml
dev_clickhouse_server_container=dev_clickhouse_server
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Starts the environment
	docker-compose -f ${dc_path} build
	docker-compose -f ${dc_path} up -d

down: ## Stops the environment
	docker-compose -f ${dc_path} down

restart: down up # Restart the environment

rebuild: ## Rebuilds the environment from scratch
	@/bin/echo -n "All the volumes will be deleted. You will loose data in DB. Are you sure? [y/N]: " && read answer && \
	[[ $${answer:-N} = y ]] && make destroy

destroy: ## Destroys thel environment
	docker-compose -f ${dc_path} down --volumes --remove-orphans
	docker-compose -f ${dc_path} rm -vsf

jump_clickhouse_server: ## Jump to the clickhouse_server container
	docker-compose -f ${dc_path} exec ${dev_clickhouse_server_container} bash
