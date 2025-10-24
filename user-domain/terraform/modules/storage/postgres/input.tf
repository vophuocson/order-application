variable "engine_type" {
  type = string
  description = "the engine runing db"
}

variable "storage_size" {
  type = number
  description = "size of storage db"
}

variable "instance_class" {
  type = string
  description = "the class  that instance using"
}

variable "skip_final_snapshot" {
  type = bool
  default = true
}

data "aws_secretsmanager_secret_version" "creds" {
  secret_id = "db-creds"
}

locals {
  db_cres = jsondecode(data.aws_secretsmanager_secret_version.creds.secret_string)
}