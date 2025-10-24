# Production Environment Configuration

environment = "prod"

# Network
vpc_cidr = "10.0.0.0/16"

# ECS - Production-ready resources
ecs_task_cpu      = 512
ecs_task_memory   = 1024
ecs_desired_count = 2
ecs_min_capacity  = 2
ecs_max_capacity  = 20

# RDS - Production configuration
rds_instance_class          = "db.t3.small"
rds_allocated_storage       = 50
rds_max_allocated_storage   = 200
rds_multi_az                = true
rds_backup_retention_period = 30
rds_skip_final_snapshot     = false

# Logging
log_retention_days = 30
enable_flow_logs   = true

# SSL/TLS - Add your ACM certificate ARN
# certificate_arn = "arn:aws:acm:us-east-1:123456789012:certificate/..."


