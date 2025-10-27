variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}


variable "database_subnet_group_name" {
  description = "Database subnet group name"
  type        = string
}

variable "allowed_security_group_ids" {
  description = "Security group IDs allowed to access RDS"
  type        = list(string)
  default     = []
}

variable "instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "Maximum allocated storage for autoscaling in GB"
  type        = number
  default     = 100
}

variable "engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16.3"
}

variable "database_name" {
  description = "Database name"
  type        = string
}

variable "master_username" {
  description = "Master username"
  type        = string
  default     = "postgres"
  sensitive   = true
}

variable "backup_retention_period" {
  description = "Backup retention period in days"
  type        = number
  default     = 7
}

variable "multi_az" {
  description = "Enable Multi-AZ deployment"
  type        = bool
  default     = false
}

variable "skip_final_snapshot" {
  description = "Skip final snapshot on deletion"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}

data "terraform_remote_state" "" {
  backend = "s3"
  config = {
    bucket = var.db_remote_state_bucket
    key = var.db_remote_state_key
    region = "ap-southeast-1"
  }
}

variable "bucket" {
  description = "The bucket name that stores the state file"
  type = string
}

variable "network_state_key" {
  description = "the key name that stores the state file"
  type = string
}

variable "region" {
  description = "the key name that stores the state file"
  type = string
}

data "terraform_remote_state" "vcp" {
  backend = "s3"
  config = {
    bucket = var.bucket
    key    = var.network_state_key
    region = var.region
  }
}