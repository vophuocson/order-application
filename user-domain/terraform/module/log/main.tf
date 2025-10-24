locals {
  name = "${var.project_name}-${var.environment}"
}

resource "aws_cloudwatch_log_group" "app" {
  name = "ecs/${local.name}"
  retention_in_days = var.log_retention_days
  tags = var.tags
}
