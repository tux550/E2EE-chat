variable "region" {
  default = "us-east-1"
}

variable "app_name" {
  default = "my-docker-app"
}

variable "ecr_image_url" {
  description = "URL of the Docker image in ECR"
  type        = string
}

variable "db_user" {
  description = "admin"
  type        = string
}

variable "db_password" {
  description = "password"
  type        = string
}