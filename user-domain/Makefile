ENV ?= development
API_DOC_DIR := ./schemas
API_DOC := $(API_DOC_DIR)/build-bundle/$(TAG)/index.yaml
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

migrate:
	go run ./cmd/migrate/main.go


dbgen:
	go run ./cmd/dbgen/main.go

apigen: 
	$(call _apigen,$(TAG))

bundle_schemas:
	swagger-cli bundle schemas/$(TAG)/index.yaml --outfile ./schemas/build-bundle/$(TAG)/index.yaml --type yaml

build:
	docker compose build

up:
	$(call _up)

down:
	docker compose down
