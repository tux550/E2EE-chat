terraform {

  required_providers {

    aws = {

      source = "hashicorp/aws"

      version = "~> 5.0"

    }

  }

}

// AWS
provider "aws" {
  region = var.region
}


// ECS Roles
resource "aws_iam_role" "ecs_task_execution" { 
  name = "${var.app_name}-ecs-task-execution-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow",
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      },
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_ecr" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_s3" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_ddb" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_docdb" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDocDBFullAccess"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_apigw" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess"
}

// Add log 
resource "aws_iam_role_policy_attachment" "ecs_task_execution_logs" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"
}

// Create log group
resource "aws_cloudwatch_log_group" "main" {
  name              = "/ecs/${var.app_name}-logs"
  retention_in_days = 7
}

// ECS
resource "aws_ecs_cluster" "main" {
  name = "${var.app_name}-cluster"
}

resource "aws_ecs_task_definition" "app" {
  family                   = "${var.app_name}-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024
  memory                   = 2048
  task_role_arn      = aws_iam_role.ecs_task_execution.arn
  execution_role_arn = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([{
    name      = "${var.app_name}-container"
    image     = var.ecr_image_url
    //var.ecr_image_url
    cpu       = 1024
    memory    = 2048
    essential = true
    portMappings = [{
      containerPort = 8080
      hostPort      = 8080
    }]
    // Log config
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = "/ecs/${var.app_name}-logs"
        "awslogs-region"        = var.region
        "awslogs-stream-prefix" = "ecs"
      }
    }

    environment = [
      // APP env
      {
        name  = "MONGO_DB"
        value = var.db_database
      },
      {
        name  = "MONGO_URI"
        value = "mongodb://${var.db_user}:${var.db_password}@${aws_docdb_cluster.main.endpoint}/?ssl=true&retryWrites=false"
      },
      {
        name  = "MONGO_CLIENT_COL"
        value = var.db_client_collection
      },
      {
        name  = "MONGO_MESSAGE_COL"
        value = var.db_message_collection
      },
      {
        name  = "DDB_TABLE_CONN"
        value = var.ddb_table_connections
      },
      {
        name  = "API_ENDPOINT",
        value = var.ws_api_gateway_endpoint
      }
    ]
  }])
}

// Load balancer
resource "aws_lb" "main" {
  name               = "${var.app_name}-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.ecs.id]
  subnets            = aws_subnet.public[*].id
}
resource "aws_lb_listener" "main" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.main.arn
  }
}
resource "aws_lb_target_group" "main" {
  name     = "${var.app_name}-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  target_type = "ip"
  health_check {
    path                = "/"
  }
}



resource "aws_ecs_service" "app" {
  name            = "${var.app_name}-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  launch_type     = "FARGATE"
  
  desired_count = 1

  force_new_deployment = true

  network_configuration {
    subnets         = aws_subnet.public[*].id
    security_groups = [aws_security_group.ecs.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.main.arn
    container_name   = "${var.app_name}-container"
    container_port   = 8080
  }
}



