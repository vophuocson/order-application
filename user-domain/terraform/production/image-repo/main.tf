module "image_repo" {
  source = "../../modules/image-repo"
  project_name = "user-api"
  environment = "production"
}

terraform {
  backend "s3" {
    bucket       = "production-terraform-up-and-running-state"
    key          = "image_repo/terraform.tfstate"
    encrypt      = true
    use_lockfile = true
    region       = "ap-southeast-1"
  }
}