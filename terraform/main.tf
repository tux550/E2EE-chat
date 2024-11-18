terraform {

  required_providers {

    aws = {

      source  = "hashicorp/aws"

      version = "~> 5.0"

    }

  }

}



// AWS
provider "aws" {
  region = var.region
}

// Create VPC and subnets
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "public" {
  count = 2
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  map_public_ip_on_launch = true
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

// Security groups
resource "aws_security_group" "ecs" {
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "documentdb" {
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 27017
    to_port     = 27017
    protocol    = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

// ECS
resource "aws_ecs_cluster" "main" {
  name = "${var.app_name}-cluster"
}


resource "aws_ecs_task_definition" "app" {
  family                   = var.app_name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([{
    name      = var.app_name
    image     = var.ecr_image_url
    cpu       = 256
    memory    = 512
    essential = true
    portMappings = [{
      containerPort = 8080
      hostPort      = 80
    }]
    environment = [
      {
        name  = "DB_URI"
        value = aws_docdb_cluster.main.endpoint
      },
      {
        name  = "DB_USER"
        value = var.db_user
      },
      {
        name  = "DB_PASSWORD"
        value = var.db_password
      }
    ]
  }])
}


resource "aws_ecs_service" "app" {
  name            = var.app_name
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = aws_subnet.public[*].id
    security_groups = [aws_security_group.ecs.id]
  }
}



// Create documentdb (Free tier)

resource "aws_docdb_cluster" "main" {
  cluster_identifier = "${var.app_name}-docdb"
  engine = "docdb"
  master_username = var.db_user
  master_password = var.db_password
  vpc_security_group_ids  = [aws_security_group.documentdb.id]
}

resource "aws_docdb_cluster_instance" "main" {
  count = 1
  identifier = "${var.app_name}-docdb-instance-${count.index}"
  cluster_identifier = aws_docdb_cluster.main.id
  instance_class = "db.t3.medium"
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
}