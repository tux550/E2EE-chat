// Create documentdb (Free tier)
resource "aws_db_subnet_group" "main" {
  name       = "${var.app_name}-docdb-subnet-group"
  subnet_ids = aws_subnet.public[*].id
}

resource "aws_docdb_cluster" "main" {
  cluster_identifier     = "${var.app_name}-docdb"
  engine                 = "docdb"
  master_username        = var.db_user
  master_password        = var.db_password
  vpc_security_group_ids = [aws_security_group.documentdb.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  skip_final_snapshot    = true
}

resource "aws_docdb_cluster_instance" "main" {
  count              = 3
  identifier         = "${var.app_name}-docdb-instance-${count.index}"
  cluster_identifier = aws_docdb_cluster.main.id
  instance_class     = "db.t3.medium"
}