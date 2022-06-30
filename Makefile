.PHONY: help
.DEFAULT_GOAL := help
SHELL=bash

dc_path=./docker-compose.yaml
dev_clickhouse_server_container=dev_clickhouse_server
dev_front_basicrum_go_container=dev_front_basicrum_go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Starts a local environment
	docker-compose -f ${dc_path} build
	docker-compose -f ${dc_path} up -d

down: ## Stops a local environment
	docker-compose -f ${dc_path} down

restart: down up # Restart environment

rebuild: ## Rebuild local environment from scratch
	@/bin/echo -n "All the volumes will be deleted. You will loose data in DB. Are you sure? [y/N]: " && read answer && \
	[[ $${answer:-N} = y ]] && make destroy

destroy: ## Destroy local environment
	docker-compose -f ${dc_path} down --volumes --remove-orphans
	docker-compose -f ${dc_path} rm -vsf

jump_front_basicrum: ## Jump to the front_basicrum container
	docker-compose -f ${dc_path} exec ${dev_front_basicrum_go_container} bash

logs_front_basicrum: ## Log messages of the front_basicrum container
	docker-compose -f ${dc_path} logs ${front_basicrum_go_container}

restart_front_basicrum: ## Restart the front_basicrum container
	docker-compose -f ${dc_path} restart ${front_basicrum_go_container}

jump_clickhouse_server: ## Jump to the clickhouse_server container
	docker-compose -f ${dc_path} exec ${dev_clickhouse_server_container} bash
