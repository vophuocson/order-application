output "alb_arn" {
  description = "ARN of the Application Load Balancer"
  value       = aws_alb.main.arn
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_alb.main.dns_name
}

output "alb_zone_id" {
  description = "Zone ID of the Application Load Balancer"
  value       = aws_alb.main.zone_id
}

output "target_group_arn" {
  description = "ARN of the Target Group"
  value       = aws_lb_target_group.app.arn
}

output "alb_security_group_id" {
  description = "Security group ID of the ALB"
  value       = aws_security_group.alb.id
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
