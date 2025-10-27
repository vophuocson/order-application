module "rds" {
  source = "../../modules/database"
  vpc_id = ""
  database_name = ""
  database_subnet_group_name = ""
  project_name = ""
  environment = ""
}


terraform {
  backend "s3" {
    bucket       = "production-terraform-up-and-running-state"
    key          = "storage/terraform.tfstate"
    encrypt      = true
    use_lockfile = true
    region       = "ap-southeast-1"
  }
}
