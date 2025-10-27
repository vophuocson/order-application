variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "container_image" {
  description = "Docker container image"
  type        = string
}

variable "container_port" {
  type        = number
  description = "Container port"
  default     = 8080
}

variable "task_cpu" {
  description = "Task CPU units"
  type        = number
  default     = 512
}

variable "task_memory" {
  description = "Task memory in MB"
  type        = number
  default     = 1024
}

variable "desired_count" {
  description = "Desired number of tasks"
  type        = number
  default     = 2
}

variable "min_capacity" {
  description = "Minimum capacity for autoscaling"
  type        = number
  default     = 1
}

variable "max_capacity" {
  description = "Maximum capacity for autoscaling"
  type        = number
  default     = 10
}


variable "health_check_path" {
  description = "Health check path"
  type        = string
  default     = "/health"
}

variable "environment_variables" {
  description = "Environment variables for container"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "secrets" {
  description = "Secrets from AWS Secrets Manager"
  type = list(object({
    name      = string
    valueFrom = string
  }))
  default = []
}

variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}

variable "cloudwatch_log_group" {
  type        = string
  description = "aws cloudwatch loggroup"
}

variable "lb_target_group" {
  type        = string
  description = "aws lb target group"
}

variable "ecs_cluster_name" {
  type        = string
  description = "ECS Cluster name"
}

variable "bucket" {
  description = "The bucket name that stores the state file"
  type = string
}

variable "vpc_state_key" {
  description = "the key name that stores the state file"
  type = string
}

variable "region" {
  description = "the key name that stores the state file"
  type = string
}

data "terraform_remote_state" "vpc" {
  backend = "s3"
  config = {
    bucket = var.bucket
    key    = var.vpc_state_key
    region = var.region
  }
}