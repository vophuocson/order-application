# ============================================================================
# Production Environment Configuration
# ============================================================================

# General
environment  = "production"
project_name = "user-api"
region       = "ap-southeast-1"
owner        = "DevOps Team"

# Backend
backend_s3_bucket = "production-terraform-up-and-running-state"

# Authentication - GitHub OIDC
allowed_repos_branches = [{
  org    = "vophuocson"
  repo   = "order-application"
  branch = "main"
}]

# Networking
vpc_cidr           = "10.0.0.0/16"
availability_zones = ["apse1-az1", "apse1-az2"]

# Database
database_name               = "order"
db_creds_secret_name        = "production/db-creds"
rds_instance_class          = "db.t3.small"
rds_allocated_storage       = 50
rds_max_allocated_storage   = 200
rds_multi_az                = true
rds_backup_retention_period = 30
rds_skip_final_snapshot     = false

# ECS - Production-ready resources
ecs_cluster_name = "order-app"
ecs_task_cpu     = 512
ecs_task_memory  = 1024
ecs_desired_count = 2
ecs_min_capacity  = 2
ecs_max_capacity  = 20

# Container Image (leave empty to use latest from ECR)
# container_image = "123456789012.dkr.ecr.ap-southeast-1.amazonaws.com/user-api:v1.0.0"
container_image = ""

# Logging
log_retention_days = 30
enable_flow_logs   = true

# SSL/TLS - Add your ACM certificate ARN for HTTPS
# certificate_arn = "arn:aws:acm:ap-southeast-1:123456789012:certificate/..."
certificate_arn = ""

