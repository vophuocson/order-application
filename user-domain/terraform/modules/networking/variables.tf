variable "project_name" {
  type        = string
  description = "Project name"
}

variable "environment" {
  type        = string
  description = "Environment name"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR blocks"
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
}

variable "enable_flow_logs" {
  description = "Enable VCP flow logs"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Additional tag"
  type        = map(string)
  default = {
  }
}