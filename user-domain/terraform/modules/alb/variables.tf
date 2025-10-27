variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}


variable "tags" {
  description = "Additional tags"
  type        = map(string)
  default     = {}
}

variable "container_port" {
  type        = number
  description = "Container port"
  default     = 8080
}

variable "certificate_arn" {
  description = "ACM certificate ARN for HTTPS"
  type        = string
  default     = ""
}


variable "health_check_path" {
  description = "Health check path"
  type        = string
  default     = "/health"
}
