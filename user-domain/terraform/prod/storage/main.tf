module "rds" {
  source = "../../modules/database"
  database_name = ""
  project_name = ""
  environment = ""
  bucket = "production-terraform-up-and-running-state"
  network_state_key = "storage/terraform.tfstate"
  region = "ap-southeast-1"
  db-creds = "db-creds"
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
