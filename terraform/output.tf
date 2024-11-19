output "ecs_cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "ecs_service_name" {
  value = aws_ecs_service.app.name
}

output "documentdb_uri" {
  value = "mongodb://${var.db_user}:${var.db_password}@${aws_docdb_cluster.main.endpoint}?ssl=true&retryWrites=false"
}

output "alb_public_url" {
  value = aws_lb.main.dns_name
}