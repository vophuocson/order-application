# Development Environment Configuration

environment = "dev"

# Network
vpc_cidr = "10.0.0.0/16"

# ECS - Minimal resources for development
ecs_task_cpu      = 256
ecs_task_memory   = 512
ecs_desired_count = 1
ecs_min_capacity  = 1
ecs_max_capacity  = 4

# RDS - Cost-optimized for development
rds_instance_class          = "db.t3.micro"
rds_allocated_storage       = 20
rds_max_allocated_storage   = 50
rds_multi_az                = false
rds_backup_retention_period = 1
rds_skip_final_snapshot     = true

# Logging
log_retention_days = 3
enable_flow_logs   = false


