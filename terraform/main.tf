terraform {

  required_providers {

    aws = {

      source  = "hashicorp/aws"

      version = "~> 5.0"

    }

  }

}

provider "aws" {
  region = "us-east-2"
}

// Sample EC2
resource "aws_instance" "first_ec2_instance" {
  ami           = "ami-0942ecd5d85baa812"
  instance_type = "t2.micro"
}

// Create documentdb (Free tier)
/*
resource "aws_docdb_cluster" "docdb_demo" {
  cluster_identifier = "my-docdb-cluster"
  engine = "docdb"
  master_username = "admin"
  master_password = "password"
}

resource "aws_docdb_cluster_instance" "docdb_instance" {
  cluster_identifier = aws_docdb_cluster.docdb_demo.id
  instance_class = "db.t3.medium"
  identifier = aws_docdb_cluster.docdb_demo.id
}
*/