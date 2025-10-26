terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

  # Uncomment and configure for remote state
  # backend "s3" {
  #   bucket         = "your-terraform-state-bucket"
  #   key            = "user-api/${var.environment}/terraform.tfstate"
  #   region         = "us-east-1"
  #   encrypt        = true
  #   dynamodb_table = "terraform-lock"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

# Data sources
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  name = "${var.project_name}-${var.environment}"
  azs  = slice(data.aws_availability_zones.available.names, 0, 3)

  common_tags = {
    Name        = local.name
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "Terraform"
  }
}

# Networking Module
module "networking" {
  source = "./modules/networking"

  project_name       = var.project_name
  environment        = var.environment
  vpc_cidr           = var.vpc_cidr
  availability_zones = local.azs
  enable_flow_logs   = var.enable_flow_logs

  tags = local.common_tags
}

# Database Module
module "database" {
  source = "./modules/database"

  project_name               = var.project_name
  environment                = var.environment
  vpc_id                     = module.networking.vpc_id
  database_subnet_group_name = module.networking.database_subnet_group_name

  # Allow access from ECS security group
  allowed_security_group_ids = [module.ecs.ecs_security_group_id]

  instance_class          = var.rds_instance_class
  allocated_storage       = var.rds_allocated_storage
  max_allocated_storage   = var.rds_max_allocated_storage
  engine_version          = var.rds_engine_version
  database_name           = var.rds_database_name
  master_username         = var.rds_username
  backup_retention_period = var.rds_backup_retention_period
  multi_az                = var.rds_multi_az
  skip_final_snapshot     = var.rds_skip_final_snapshot

  tags = local.common_tags
}

# ECR Module
module "image_repo" {
  source = "./modules/image-repo"

  project_name = var.project_name
  environment  = var.environment

  tags = local.common_tags
}

# CloudWatch Log Module
module "log" {
  source = "./modules/log"

  project_name       = var.project_name
  environment        = var.environment
  log_retention_days = var.log_retention_days

  tags = local.common_tags
}

# ALB Module
module "alb" {
  source = "./modules/alb"

  project_name      = var.project_name
  environment       = var.environment
  vpc_id            = module.networking.vpc_id
  public_subnet_ids = module.networking.public_subnet_ids
  container_port    = var.container_port
  health_check_path = var.health_check_path
  certificate_arn   = var.certificate_arn

  tags = local.common_tags
}

# ECS Module
module "ecs" {
  source = "./modules/ecs"

  project_name       = var.project_name
  environment        = var.environment
  vpc_id             = module.networking.vpc_id
  private_subnet_ids = module.networking.private_subnet_ids

  container_image       = var.container_image
  container_port        = var.container_port
  task_cpu              = var.ecs_task_cpu
  task_memory           = var.ecs_task_memory
  desired_count         = var.ecs_desired_count
  min_capacity          = var.ecs_min_capacity
  max_capacity          = var.ecs_max_capacity
  health_check_path     = var.health_check_path
  cloudwatch_log_group  = module.log.log_group_name
  lb_target_group       = module.alb.target_group_arn
  ecs_cluster_name      = module.alb.ecs_cluster_name
  ecs_security_group_id = module.alb.ecs_security_group_id

  environment_variables = [
    {
      name  = "APP_ENV"
      value = var.environment
    },
    {
      name  = "API_PORT"
      value = tostring(var.container_port)
    },
    {
      name  = "SECRET_POSTGRES_HOSTNAME"
      value = module.database.db_instance_address
    },
    {
      name  = "SECRET_POSTGRES_PORT"
      value = tostring(module.database.db_instance_port)
    },
    {
      name  = "SECRET_POSTGRES_DATABASE"
      value = module.database.db_name
    },
    {
      name  = "SECRET_POSTGRES_USER"
      value = module.database.db_master_username
    },
    {
      name  = "SECRET_POSTGRES_SSL_MODE"
      value = "require"
    }
  ]

  secrets = [
    {
      name      = "SECRET_POSTGRES_PASSWORD"
      valueFrom = "${module.database.db_secret_arn}:password::"
    }
  ]

  alb_security_id = module.alb.alb_security_group_id

  tags = local.common_tags
}

