module "s3" {
  source = "../../modules/bucket"
  bucket_name = "production-terraform-up-and-running-state"
}