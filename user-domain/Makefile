ENV ?= development
API_DOC_DIR := ./schemas
API_DOC := $(API_DOC_DIR)/build-bundle/$(TAG).yaml
OUTPUT_DIR := ./pkg/handler
CONFIG_FILE:=$(API_DOC_DIR)/oapi-codegen.yaml

define _apigen
	oapi-codegen -config=$(CONFIG_FILE) -include-tags="$(1)" $(API_DOC) > $(OUTPUT_DIR)/$(1).gen.go
endef

define _up
	if [ "$(ENV)" = "production" ]; then \
		docker compose up -d $(SERVICE); \
	else \
		docker compose --env-file .env up -d --build $(SERVICE); \
	fi
endef

apigen: 
	$(call _apigen,$(TAG))

dbgen:
	go run ./cmd/dbgen/main.go

migrate:
	go run ./cmd/migrate/main.go

build:
	docker compose build

up:
	$(call _up)

down:
	docker compose down
