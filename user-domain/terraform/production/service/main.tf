module "alb" {
  source = "../../modules/alb"
  bucket = "production-terraform-up-and-running-state"
  project_name = "user-api"
  region       = "ap-southeast-1"
  environment = "production"
  network_state_key =  "networking/terraform.tfstate"
}

module "ecs" {
  source = "../../modules/ecs"
  bucket = "production-terraform-up-and-running-state"
  ecs_cluster_name = "order-app"
  container_image = ""
  project_name = "user-api"
  cloudwatch_log_group = ""
  environment = "production"
  vpc_state_key = "networking/terraform.tfstate"
  # network_state_key = "/terraform.tfstate"
  region = "ap-southeast-1"
  alb_target_group = module.alb.target_group_arn
}

terraform {
  backend "s3" {
    bucket       = "production-terraform-up-and-running-state"
    key          = "service/terraform.tfstate"
    encrypt      = true
    use_lockfile = true
    region       = "ap-southeast-1"
  }
}