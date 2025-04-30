SHELL := /usr/bin/env bash

CLUSTER_SCRIPT := ./generate-cluster-id.sh
COMPOSE := docker compose --env-file .env

.PHONY: up down logs force

force:

up: .env
	$(COMPOSE) up -d

.env: force
	@$(CLUSTER_SCRIPT)

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f