module "vcp" {
  project_name = "svp-order"
  source = "../../modules/networking"
  availability_zones = ["apse1-az1"]
  environment = "production"
}
