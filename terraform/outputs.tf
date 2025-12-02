output "ecr_repository_url" {
  description = "URL of the ECR repository"
  value       = aws_ecr_repository.main.repository_url
}

output "ecr_repository_arn" {
  description = "ARN of the ECR repository"
  value       = aws_ecr_repository.main.arn
}

output "github_actions_role_arn" {
  description = "ARN of the IAM role for GitHub Actions (set this as AWS_ROLE_ARN secret)"
  value       = aws_iam_role.github_actions.arn
}

output "oidc_provider_arn" {
  description = "ARN of the GitHub OIDC provider (existing or newly created)"
  value       = local.oidc_provider_arn
}

output "oidc_provider_created" {
  description = "Whether the OIDC provider was created by this module"
  value       = var.create_oidc_provider
}
