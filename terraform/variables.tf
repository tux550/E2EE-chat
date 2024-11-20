variable "region" {
  default = "us-east-2"
}

variable "availability_zone" {
  description = "Availability zone"
  default = ["us-east-2a",
  "us-east-2b"]
  type = list(string)

}


variable "app_name" {
  default = "my-docker-app"
}



variable "ecr_image_url" {
  description = "URL of the Docker image in ECR"
  type        = string
}

variable "ecr_demo_image_url" {
  description = "URL of the Docker image in ECR for demo"
  type        = string
  default = ""
  
}

variable "db_user" {
  description = "Database user"
  type        = string
  default     = "myuser"
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
  default     = "e2ee_chat_connections"
}
