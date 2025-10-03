API_DOC_DIR := ./schemas
API_DOC := $(API_DOC_DIR)/build-bundle/$(TAG).yaml
OUTPUT_DIR := ./pkg/handler
CONFIG_FILE:=$(API_DOC_DIR)/oapi-codegen.yaml

define _apigen
	oapi-codegen -config=$(CONFIG_FILE) -include-tags="$(1)" $(API_DOC) > $(OUTPUT_DIR)/$(1).gen.go
endef

apigen: 
	$(call _apigen,$(TAG))
