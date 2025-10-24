variable "project_name" {
  description = "Project name"
  type = string
}

variable "environment" {
  description = "Environment name"
  type = string
}


variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs for ALB"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for ECS tasks"
  type        = list(string)
}

variable "container_image" {
  description = "Docker container image"
  type        = string
}

variable "container_port" {
  type = number
  description = "Container port"
  default = 8080
}

variable "task_cpu" {
  description = "Task CPU units"
  type = number
  default = 512
}

variable "task_memory" {
  description = "Task memory in MB"
  type = number
  default = 1024
}

variable "desired_cou t" {
  description = "Desired number of tasks"
  type = number
  default = 2
}

variable "min_capacity" {
  description = "Minimum capacity for autoscaling"
  type = number
  default = 1
}

variable "max_capacity" {
  description = "Maximum capacity for autoscaling"
  type = number
  default = 10
}


variable "health_check_path" {
  description = "Health check path"
  type        = string
  default     = "/health"
}

variable "environment_variables" {
  description = "Environment variables for container"
  type = list(object({
    name = string
    value = string
  }))
  default = []
}

variable "secret" {
  description = "Secrets from AWS Secrets Manager"
  type = list(object({
    name = string
    valueFrom = string
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


variable "desired_count" {
  description = "Desired number of tasks"
  type        = number
  default     = 2
}

variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}

variable "cloudwatch_log_group" {
  type = string
  description = "aws cloudwatch loggroup"
}

variable "lb_target_group" {
  type = string
  description = "aws lb target group"
}