locals {
  name = "${var.project_name}-${var.environment}"
}

# Security Group for RDS
resource "aws_security_group" "rds" {
  name_prefix = "${local.name}-rds-"
  description = "Security group for RDS database"
  vpc_id      = data.terraform_remote_state.vcp.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = var.allowed_security_group_ids
    description     = "PostgreSQL access from ECS tasks"
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow HTTPS outbound for monitoring or updates"
  }

  tags = merge(var.tags, {
    Name = "${local.name}-rds-sg"
  })

  lifecycle {
    create_before_destroy = true
  }
}

# RDS Instance
resource "aws_db_instance" "main" {
  identifier_prefix = "${var.project_name}-${var.environment}-"

  engine                = "postgres"
  engine_version        = var.engine_version
  instance_class        = var.instance_class
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_encrypted     = true

  db_name  = var.database_name
  username = local.db_cres.username
  password = local.db_cres.password

  db_subnet_group_name   = var.database_subnet_group_name
  vpc_security_group_ids = [aws_security_group.rds.id]

  multi_az                = var.multi_az
  publicly_accessible     = false
  backup_retention_period = var.backup_retention_period
  backup_window           = "03:00-04:00"
  maintenance_window      = "mon:04:00-mon:05:00"

  skip_final_snapshot       = var.skip_final_snapshot
  final_snapshot_identifier = var.skip_final_snapshot ? null : "${local.name}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  tags = var.tags
}

