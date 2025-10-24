variable "project_name" {
  description = "Project name"
  type = string
}

variable "environment" {
  description = "Environment name"
  type = string
}


variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs for ALB"
  type        = list(string)
}

variable "container_port" {
  type = number
  description = "Container port"
  default = 8080
}