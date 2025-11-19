output "oidc_provider_arn" {
  description = "ARN of the GitHub OIDC provider"
  value       = aws_iam_openid_connect_provider.github_action.arn
}

output "auth_role_arn" {
  description = "ARN of the GitHub Actions authentication role"
  value       = aws_iam_role.auth_role.arn
}

output "auth_role_name" {
  description = "Name of the GitHub Actions authentication role"
  value       = aws_iam_role.auth_role.name
}



