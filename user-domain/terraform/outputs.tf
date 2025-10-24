# VPC Outputs
output "vpc_id" {
  description = "VPC ID"
  value       = module.networking.vpc_id
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = module.networking.public_subnet_ids
}

output "private_subnet_ids" {
  description = "Private subnet IDs"
  value       = module.networking.private_subnet_ids
}

# Database Outputs
output "db_endpoint" {
  description = "RDS database endpoint"
  value       = module.database.db_instance_endpoint
}

output "db_name" {
  description = "Database name"
  value       = module.database.db_name
}

output "db_secret_arn" {
  description = "ARN of Secrets Manager secret containing DB credentials"
  value       = module.database.db_secret_arn
  sensitive   = true
}

# ECR Outputs
output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = module.image_repo.repository_url
}

# ALB Outputs
output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.alb.alb_dns_name
}

output "alb_zone_id" {
  description = "Zone ID of the Application Load Balancer"
  value       = module.alb.alb_zone_id
}

# ECS Outputs
output "ecs_cluster_name" {
  description = "ECS Cluster name"
  value       = module.alb.ecs_cluster_name
}

# CloudWatch Outputs
output "log_group_name" {
  description = "CloudWatch Log Group name"
  value       = module.log.log_group_name
}

