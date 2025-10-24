.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Variables
ENV ?= development
AWS_REGION ?= us-east-1
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
PROJECT_NAME ?= user-api
ECR_REPO ?= $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(PROJECT_NAME)-$(ENV)
IMAGE_TAG ?= latest

API_DOC_DIR := ./schemas
API_DOC := $(API_DOC_DIR)/build-bundle/$(TAG)/index.yaml
OUTPUT_DIR := ./infrastructure/http/handler
CONFIG_FILE := $(API_DOC_DIR)/oapi-codegen.yaml

# ==================== Local Development ====================

.PHONY: install
install: ## Install dependencies
	go mod download
	go mod verify

.PHONY: run
run: ## Run the application locally
	go run ./cmd/api/main.go

.PHONY: test
test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: coverage
coverage: test ## Show test coverage
	go tool cover -html=coverage.out

.PHONY: lint
lint: ## Run linter
	golangci-lint run ./...

# ==================== Code Generation ====================

.PHONY: apigen
apigen: ## Generate API code from OpenAPI spec (usage: make apigen TAG=user)
	oapi-codegen -config=$(CONFIG_FILE) -include-tags="$(TAG)" $(API_DOC) > $(OUTPUT_DIR)/$(TAG).gen.go

.PHONY: bundle_schemas
bundle_schemas: ## Bundle OpenAPI schemas (usage: make bundle_schemas TAG=user)
	swagger-cli bundle schemas/$(TAG)/index.yaml --outfile ./schemas/build-bundle/$(TAG)/index.yaml --type yaml

.PHONY: dbgen
dbgen: ## Generate database code
	go run ./cmd/dbgen/main.go

.PHONY: migrate
migrate: ## Run database migrations
	go run ./cmd/migrate/main.go

.PHONY: mock_domain
mock_domain: ## Generate domain mocks
	mockery --all --dir=./internal/domain --disable-func-mocks=true --keeptree=true --recursive=true --output=./internal/domain/mocks --case=underscore --outpkg=domainmock

.PHONY: mock_application
mock_application: ## Generate application mocks
	mockery --all --dir=./internal/application --disable-func-mocks=true --keeptree=true --recursive=true --output=./internal/application/mocks --case=underscore --outpkg=applicationmock

.PHONY: mock_infrastructure
mock_infrastructure: ## Generate infrastructure mocks
	mockery --all --dir=./infrastructure --disable-func-mocks=true --keeptree=true --recursive=true --output=./infrastructure/mocks --case=underscore --outpkg=infrastructure_mock

# ==================== Docker Commands ====================

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(PROJECT_NAME):$(IMAGE_TAG) .

.PHONY: docker-up
docker-up: ## Start Docker Compose services
	docker compose up -d --build

.PHONY: docker-down
docker-down: ## Stop Docker Compose services
	docker compose down

.PHONY: docker-logs
docker-logs: ## View Docker Compose logs
	docker compose logs -f

.PHONY: docker-clean
docker-clean: ## Clean Docker images and volumes
	docker compose down -v
	docker system prune -f

# ==================== AWS Commands ====================

.PHONY: aws-login
aws-login: ## Login to AWS ECR
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

.PHONY: ecr-create
ecr-create: ## Create ECR repository
	aws ecr describe-repositories --repository-names $(PROJECT_NAME)-$(ENV) --region $(AWS_REGION) || \
	aws ecr create-repository --repository-name $(PROJECT_NAME)-$(ENV) --region $(AWS_REGION)

.PHONY: docker-push
docker-push: aws-login ## Build and push Docker image to ECR
	docker build --platform linux/amd64 -t $(PROJECT_NAME):$(IMAGE_TAG) .
	docker tag $(PROJECT_NAME):$(IMAGE_TAG) $(ECR_REPO):$(IMAGE_TAG)
	docker push $(ECR_REPO):$(IMAGE_TAG)

.PHONY: docker-push-versioned
docker-push-versioned: ## Push versioned Docker image (usage: make docker-push-versioned VERSION=1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	docker build --platform linux/amd64 -t $(PROJECT_NAME):$(VERSION) .
	docker tag $(PROJECT_NAME):$(VERSION) $(ECR_REPO):$(VERSION)
	docker tag $(PROJECT_NAME):$(VERSION) $(ECR_REPO):latest
	docker push $(ECR_REPO):$(VERSION)
	docker push $(ECR_REPO):latest

# ==================== Terraform Commands ====================

.PHONY: tf-init
tf-init: ## Initialize Terraform
	cd terraform && terraform init

.PHONY: tf-plan
tf-plan: ## Plan Terraform changes
	cd terraform && terraform plan -var="environment=$(ENV)"

.PHONY: tf-apply
tf-apply: ## Apply Terraform changes
	cd terraform && terraform apply -var="environment=$(ENV)"

.PHONY: tf-destroy
tf-destroy: ## Destroy Terraform resources
	cd terraform && terraform destroy -var="environment=$(ENV)"

.PHONY: tf-output
tf-output: ## Show Terraform outputs
	cd terraform && terraform output

.PHONY: tf-fmt
tf-fmt: ## Format Terraform files
	cd terraform && terraform fmt -recursive

.PHONY: tf-validate
tf-validate: ## Validate Terraform configuration
	cd terraform && terraform validate

# ==================== ECS Commands ====================

.PHONY: ecs-deploy
ecs-deploy: docker-push ## Deploy to ECS (push image and update service)
	aws ecs update-service \
		--cluster $(PROJECT_NAME)-$(ENV) \
		--service $(PROJECT_NAME)-$(ENV)-service \
		--force-new-deployment \
		--region $(AWS_REGION)

.PHONY: ecs-logs
ecs-logs: ## View ECS logs
	aws logs tail /ecs/$(PROJECT_NAME)-$(ENV) --follow --region $(AWS_REGION)

.PHONY: ecs-status
ecs-status: ## Check ECS service status
	aws ecs describe-services \
		--cluster $(PROJECT_NAME)-$(ENV) \
		--services $(PROJECT_NAME)-$(ENV)-service \
		--region $(AWS_REGION) \
		--query 'services[0].{Status:status,Running:runningCount,Desired:desiredCount,Pending:pendingCount}'

.PHONY: ecs-scale
ecs-scale: ## Scale ECS service (usage: make ecs-scale COUNT=3)
	@if [ -z "$(COUNT)" ]; then echo "COUNT is required"; exit 1; fi
	aws ecs update-service \
		--cluster $(PROJECT_NAME)-$(ENV) \
		--service $(PROJECT_NAME)-$(ENV)-service \
		--desired-count $(COUNT) \
		--region $(AWS_REGION)

.PHONY: ecs-exec
ecs-exec: ## Execute command in ECS task
	@TASK_ARN=$$(aws ecs list-tasks --cluster $(PROJECT_NAME)-$(ENV) --service-name $(PROJECT_NAME)-$(ENV)-service --region $(AWS_REGION) --query 'taskArns[0]' --output text); \
	aws ecs execute-command \
		--cluster $(PROJECT_NAME)-$(ENV) \
		--task $$TASK_ARN \
		--container $(PROJECT_NAME) \
		--interactive \
		--command "/bin/sh" \
		--region $(AWS_REGION)

# ==================== Database Commands ====================

.PHONY: db-migrate-aws
db-migrate-aws: ## Run migrations on AWS RDS
	@echo "Running migrations on AWS RDS..."
	@DB_ENDPOINT=$$(cd terraform && terraform output -raw db_instance_endpoint); \
	DB_SECRET=$$(cd terraform && terraform output -raw db_secret_arn); \
	echo "Database endpoint: $$DB_ENDPOINT"

# ==================== Deployment Pipeline ====================

.PHONY: deploy-infra
deploy-infra: tf-init tf-apply ## Deploy infrastructure
	@echo "Infrastructure deployed successfully"
	@cd terraform && terraform output

.PHONY: deploy-app
deploy-app: docker-push ecs-deploy ## Deploy application
	@echo "Application deployed successfully"
	@make ecs-status

.PHONY: deploy-all
deploy-all: deploy-infra ecr-create deploy-app ## Full deployment (infrastructure + application)
	@echo "Full deployment completed"

# ==================== Cleanup ====================

.PHONY: clean
clean: ## Clean local build artifacts
	rm -rf coverage.out
	rm -rf tmp/
	go clean

.PHONY: clean-all
clean-all: clean docker-clean ## Clean everything including Docker
	@echo "Cleanup completed"
