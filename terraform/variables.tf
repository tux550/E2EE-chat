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
  description = "Database user"
  type        = string
  default     = "admin"
}

variable "db_password" {
  description = "Database password"
  type        = string
  default     = "password"
}

variable "db_database" {
  description = "Database name"
  type        = string
  default     = "mydb"
}

variable "db_client_collection" {
  description = "Database client collection"
  type        = string
  default     = "client"
}

variable "db_message_collection" {
  description = "Database message collection"
  type        = string
  default     = "message" 
}


// TODO: Move the implementation to terraform
variable "ws_api_gateway_endpoint" {
  description = "Websocket API Gateway endpoint"
  type        = string
}

variable "ddb_table_connections" {
  description = "DynamoDB table for connections"
  type        = string
}
