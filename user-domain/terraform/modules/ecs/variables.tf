variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}


variable "vpc_id" {
  description = "VPC ID"
  type        = string
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

variable "alb_security_id" {
  type        = string
  description = "aws alb security group ID"
}

variable "ecs_cluster_name" {
  type        = string
  description = "ECS Cluster name"
}

variable "ecs_security_group_id" {
  type        = string
  description = "ECS Task Security Group ID"
}