variable "environment" {
  description = "Environment name (e.g., production, staging)"
  type        = string
  default     = "production"
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "user-api"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-1"
}

variable "owner" {
  description = "Owner of the infrastructure"
  type        = string
  default     = "DevOps Team"
}

# Backend Configuration

variable "backend_s3_bucket" {
  description = "S3 bucket name for Terraform state"
  type        = string
  default     = "production-terraform-up-and-running-state"
}

# Authentication Variables

variable "allowed_repos_branches" {
  description = "List of allowed GitHub repositories and branches for OIDC"
  type = list(object({
    org    = string
    repo   = string
    branch = string
  }))
  default = [{
    org    = "vophuocson"
    repo   = "order-application"
    branch = "main"
  }]
}

# Networking Variables

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["apse1-az1", "apse1-az2"]
}

# Database Variables

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "order"
}

variable "db_creds_secret_name" {
  description = "AWS Secrets Manager secret name for DB credentials"
  type        = string
  default     = "db-creds"
}

variable "rds_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.small"
}

variable "rds_allocated_storage" {
  description = "Initial allocated storage for RDS (GB)"
  type        = number
  default     = 50
}

variable "rds_max_allocated_storage" {
  description = "Maximum allocated storage for RDS autoscaling (GB)"
  type        = number
  default     = 200
}

variable "rds_multi_az" {
  description = "Enable Multi-AZ deployment for RDS"
  type        = bool
  default     = true
}

variable "rds_backup_retention_period" {
  description = "Number of days to retain automated backups"
  type        = number
  default     = 30
}

variable "rds_skip_final_snapshot" {
  description = "Skip final snapshot when destroying RDS instance"
  type        = bool
  default     = false
}

# ECS Variables

variable "ecs_cluster_name" {
  description = "ECS cluster name"
  type        = string
  default     = "order-app"
}

variable "container_image" {
  description = "Container image for ECS task (leave empty to use ECR latest)"
  type        = string
  default     = ""
}

variable "ecs_task_cpu" {
  description = "CPU units for ECS task"
  type        = number
  default     = 512
}

variable "ecs_task_memory" {
  description = "Memory (MB) for ECS task"
  type        = number
  default     = 1024
}

variable "ecs_desired_count" {
  description = "Desired number of ECS tasks"
  type        = number
  default     = 2
}

variable "ecs_min_capacity" {
  description = "Minimum number of ECS tasks for autoscaling"
  type        = number
  default     = 2
}

variable "ecs_max_capacity" {
  description = "Maximum number of ECS tasks for autoscaling"
  type        = number
  default     = 20
}

# Logging Variables

variable "log_retention_days" {
  description = "CloudWatch log retention period in days"
  type        = number
  default     = 30
}

variable "enable_flow_logs" {
  description = "Enable VPC flow logs"
  type        = bool
  default     = true
}

# SSL/TLS Variables

variable "certificate_arn" {
  description = "ACM certificate ARN for HTTPS (leave empty to use HTTP only)"
  type        = string
  default     = ""
}

variable "log_retention_days" {
  type = number
}