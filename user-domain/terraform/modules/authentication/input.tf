data "tls_certificate" "github" {
  url = "https://token.actions.githubusercontent.com"
}

variable "allowed_repos_branches" {
  description = "GitHub repos/branches allowed to assume the IAM role."
  type = list(object({
    org     = string
    repo    = string
    branche = string
  }))
}