module "OIDC" {
  source = "../../modules/authentication"
  allowed_repos_branches = [{
    org = "vophuocson"
    repo = "order-application"
    branche = "main"
  }]
}

terraform {
  backend "s3" {
    bucket       = "prod-terraform-up-and-running-state"
    key          = "authentication/terraform.tfstate"
    encrypt      = true
    use_lockfile = true
    region       = "ap-southeast-1"
  }
}