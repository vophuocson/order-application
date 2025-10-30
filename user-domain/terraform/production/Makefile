# ============================================================================
# Makefile for Terraform Production Environment
# ============================================================================

.PHONY: help init plan apply destroy fmt validate clean docs

# Default target
.DEFAULT_GOAL := help

# Variables
PLAN_FILE := tfplan
ENV := production

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED := \033[31m

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)Terraform Production Environment$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_YELLOW)Available targets:$(COLOR_RESET)"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':'

## init: Initialize Terraform working directory
init:
	@echo "$(COLOR_GREEN)Initializing Terraform...$(COLOR_RESET)"
	terraform init -upgrade

## validate: Validate Terraform configuration
validate:
	@echo "$(COLOR_GREEN)Validating configuration...$(COLOR_RESET)"
	terraform validate

## fmt: Format Terraform files
fmt:
	@echo "$(COLOR_GREEN)Formatting Terraform files...$(COLOR_RESET)"
	terraform fmt -recursive

## lint: Run linting checks
lint: fmt validate
	@echo "$(COLOR_GREEN)Running linting checks...$(COLOR_RESET)"
	@command -v tflint >/dev/null 2>&1 && tflint || echo "$(COLOR_YELLOW)tflint not installed, skipping...$(COLOR_RESET)"

## plan: Create Terraform execution plan
plan: validate
	@echo "$(COLOR_GREEN)Creating execution plan...$(COLOR_RESET)"
	terraform plan -out=$(PLAN_FILE)

## apply: Apply Terraform changes
apply:
	@echo "$(COLOR_YELLOW)Applying changes to $(ENV) environment...$(COLOR_RESET)"
	@if [ -f $(PLAN_FILE) ]; then \
		terraform apply $(PLAN_FILE); \
		rm -f $(PLAN_FILE); \
	else \
		echo "$(COLOR_RED)Error: No plan file found. Run 'make plan' first.$(COLOR_RESET)"; \
		exit 1; \
	fi

## apply-auto: Plan and apply in one step (use with caution!)
apply-auto: validate
	@echo "$(COLOR_YELLOW)WARNING: Auto-applying to $(ENV) environment...$(COLOR_RESET)"
	terraform apply -auto-approve

## destroy: Destroy Terraform-managed infrastructure
destroy:
	@echo "$(COLOR_RED)WARNING: This will destroy all resources in $(ENV) environment!$(COLOR_RESET)"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read confirm
	terraform destroy

## output: Show Terraform outputs
output:
	@terraform output

## output-json: Show outputs in JSON format
output-json:
	@terraform output -json

## refresh: Refresh Terraform state
refresh:
	@echo "$(COLOR_GREEN)Refreshing state...$(COLOR_RESET)"
	terraform refresh

## state-list: List resources in state
state-list:
	@terraform state list

## state-show: Show a specific resource (usage: make state-show RESOURCE=aws_instance.example)
state-show:
	@if [ -z "$(RESOURCE)" ]; then \
		echo "$(COLOR_RED)Error: RESOURCE parameter required$(COLOR_RESET)"; \
		echo "Usage: make state-show RESOURCE=aws_instance.example"; \
		exit 1; \
	fi
	@terraform state show $(RESOURCE)

## clean: Remove generated files
clean:
	@echo "$(COLOR_GREEN)Cleaning up...$(COLOR_RESET)"
	rm -f $(PLAN_FILE)
	rm -f .terraform.lock.hcl
	rm -f crash.log
	@echo "$(COLOR_YELLOW)Note: .terraform directory preserved$(COLOR_RESET)"

## clean-all: Remove all Terraform files (including .terraform)
clean-all: clean
	@echo "$(COLOR_RED)Removing .terraform directory...$(COLOR_RESET)"
	rm -rf .terraform

## docs: Generate documentation using terraform-docs
docs:
	@echo "$(COLOR_GREEN)Generating documentation...$(COLOR_RESET)"
	@command -v terraform-docs >/dev/null 2>&1 && \
		terraform-docs markdown table . > TERRAFORM_DOCS.md || \
		echo "$(COLOR_YELLOW)terraform-docs not installed. Install from: https://terraform-docs.io/$(COLOR_RESET)"

## unlock: Force unlock state (usage: make unlock LOCK_ID=xxx)
unlock:
	@if [ -z "$(LOCK_ID)" ]; then \
		echo "$(COLOR_RED)Error: LOCK_ID parameter required$(COLOR_RESET)"; \
		echo "Usage: make unlock LOCK_ID=xxx"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Force unlocking state with ID: $(LOCK_ID)$(COLOR_RESET)"
	terraform force-unlock $(LOCK_ID)

## graph: Generate dependency graph
graph:
	@echo "$(COLOR_GREEN)Generating dependency graph...$(COLOR_RESET)"
	terraform graph | dot -Tpng > graph.png
	@echo "Graph saved to graph.png"

## cost: Estimate infrastructure costs (requires Infracost)
cost:
	@echo "$(COLOR_GREEN)Estimating costs...$(COLOR_RESET)"
	@command -v infracost >/dev/null 2>&1 && \
		infracost breakdown --path . || \
		echo "$(COLOR_YELLOW)Infracost not installed. Install from: https://www.infracost.io/$(COLOR_RESET)"

## security: Run security scan (requires tfsec)
security:
	@echo "$(COLOR_GREEN)Running security scan...$(COLOR_RESET)"
	@command -v tfsec >/dev/null 2>&1 && \
		tfsec . || \
		echo "$(COLOR_YELLOW)tfsec not installed. Install from: https://github.com/aquasecurity/tfsec$(COLOR_RESET)"

## full-check: Run all checks (fmt, validate, lint, security)
full-check: fmt validate lint security
	@echo "$(COLOR_GREEN)All checks completed!$(COLOR_RESET)"

## setup: Setup backend infrastructure (S3 bucket and DynamoDB table)
setup:
	@echo "$(COLOR_GREEN)Setting up backend infrastructure...$(COLOR_RESET)"
	@bash scripts/setup-backend.sh

## workflow: Complete workflow (fmt -> validate -> plan)
workflow: fmt validate plan
	@echo "$(COLOR_GREEN)Workflow completed! Review the plan and run 'make apply' to proceed.$(COLOR_RESET)"

