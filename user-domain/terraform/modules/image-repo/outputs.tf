output "repository_url" {
  description = "ECR Repository URL"
  value       = aws_ecr_repository.app.repository_url
}

output "repository_arn" {
  description = "ECR Repository ARN"
  value       = aws_ecr_repository.app.arn
}

output "repository_name" {
  description = "ECR Repository name"
  value       = aws_ecr_repository.app.name
}

