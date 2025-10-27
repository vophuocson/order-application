module "vcp" {
  project_name = "svp-order"
  source = "../../modules/networking"
  availability_zones = ["apse1-az1"]
  environment = "production"
}


terraform {
  backend "s3" {
    bucket = "production-terraform-up-and-running-state"
    key = "networking/terraform.tfstate"
    region = "ap-southeast-1"
    encrypt = true
    use_lockfile = true
  }
}
