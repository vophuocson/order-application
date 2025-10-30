data "aws_secretsmanager_secret_version" "creds" {
  secret_id = var.db_creds
}

locals {
  name    = "${var.project_name}-${var.environment}"
  db_cres = jsondecode(data.aws_secretsmanager_secret_version.creds.secret_string)
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
  vpc_security_group_ids = var.vpc_security_group_ids

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

