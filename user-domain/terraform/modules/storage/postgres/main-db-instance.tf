resource "aws_db_instance" "db" {
  identifier_prefix   = "db"
  engine              = var.engine_type
  allocated_storage   = var.storage_size
  instance_class      = var.instance_class
  skip_final_snapshot = var.skip_final_snapshot
  username            = local.db_cres.username
  password            = local.db_cres.password
}