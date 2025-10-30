
locals {
  environment  = var.environment
  project_name = var.project_name
  region       = var.region
  
  common_tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "Terraform"
    Owner       = var.owner
  }
}

# Networking - VPC, Subnets, Security Groups

module "networking" {
  source = "../modules/networking"
  
  project_name       = local.project_name
  environment        = local.environment
  availability_zones = var.availability_zones
  vpc_cidr           = var.vpc_cidr
}

# Image Repository - ECR

module "image_repo" {
  source = "../modules/image-repo"
  
  project_name = local.project_name
  environment  = local.environment
}

# Logging - CloudWatch Log Groups

module "logging" {
  source = "../modules/log"
  
  project_name      = local.project_name
  environment       = local.environment
  log_retention_days = var.log_retention_days
}

# Database - RDS PostgreSQL

module "database" {
  source = "../modules/database"
  
  project_name     = local.project_name
  environment      = local.environment
  database_name    = var.database_name
  region           = local.region
  bucket           = var.backend_s3_bucket
  db_creds         = var.db_creds_secret_name
  

  # RDS Configuration
  instance_class          = var.rds_instance_class
  allocated_storage       = var.rds_allocated_storage
  max_allocated_storage   = var.rds_max_allocated_storage
  multi_az                = var.rds_multi_az
  backup_retention_period = var.rds_backup_retention_period
  skip_final_snapshot     = var.rds_skip_final_snapshot
  vpc_security_group_ids = [module.networking.rds_security_group]
  database_subnet_group_name = module.networking.database_subnet_group_name
}

# Application Load Balancer

module "alb" {
  source = "../modules/alb"
  
  project_name = local.project_name
  environment  = local.environment
  # Pass VPC outputs directly
  vpc_id            = module.networking.vpc_id
  public_subnet_ids = module.networking.public_subnet_ids
  
  # Optional: SSL Certificate
  certificate_arn = var.certificate_arn
  security_groups = [module.networking.alb_security_group]
}

# ECS Cluster and Service

module "ecs" {
  source = "../modules/ecs"
  
  project_name     = local.project_name
  environment      = local.environment
  region           = local.region
  bucket           = var.backend_s3_bucket
  ecs_cluster_name = var.ecs_cluster_name
  
  # Container Configuration
  repository_url = var.container_image != "" ? var.container_image : "${module.image_repo.repository_url}:latest"
  
  # CloudWatch Logs
  cloudwatch_log_group = module.logging.log_group_name
  
  # ECS Task Configuration
  task_cpu       = var.ecs_task_cpu
  task_memory    = var.ecs_task_memory
  desired_count  = var.ecs_desired_count
  min_capacity   = var.ecs_min_capacity
  max_capacity   = var.ecs_max_capacity
  
  # Networking
  private_subnet_ids = module.networking.private_subnet_ids
  security_groups    = [module.networking.ecs_security_group]
  
  # Load Balancer
  alb_target_group = module.alb.target_group_arn
  
  image_tag = ""
  
  depends_on = [
    module.alb,
    module.database,
    module.logging
  ]
}

