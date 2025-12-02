variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

variable "ecr_repository_name" {
  description = "Name of the ECR repository"
  type        = string
  default     = "go-web-server-docker-example"
}

variable "github_repository" {
  description = "GitHub repository in format 'owner/repo'"
  type        = string
}

variable "create_oidc_provider" {
  description = "Whether to create a new GitHub OIDC provider. Set to false only if one already exists in your account."
  type        = bool
  default     = true
}
