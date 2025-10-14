ENV ?= development
API_DOC_DIR := ./schemas
API_DOC := $(API_DOC_DIR)/build-bundle/$(TAG)/index.yaml
OUTPUT_DIR := ./infrastructure/http/handler
CONFIG_FILE:=$(API_DOC_DIR)/oapi-codegen.yaml

define _apigen
	oapi-codegen -config=$(CONFIG_FILE) -include-tags="$(1)" $(API_DOC) > $(OUTPUT_DIR)/$(1).gen.go
endef

define _up
	if [ "$(ENV)" = "production" ]; then \
		docker compose up -d; \
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

mock_domain:
	mockery --all --dir=./internal/domain --disable-func-mocks=true --keeptree=true --recursive=true --output=./internal/domain/mocks --case=underscore --outpkg=domain_mock

mock_application:
	mockery --all --dir=./internal/application --disable-func-mocks=true --keeptree=true --recursive=true --output=./internal/application/mocks --case=underscore --outpkg=application_mock

mock_infrastructure:
	mockery --all --dir=./infrastructure --disable-func-mocks=true --keeptree=true --recursive=true --output=./infrastructure/mocks --case=underscore --outpkg=infrastructure_mock
