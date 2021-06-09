provider "aws" {
  region = "us-east-1"
}

terraform {
  required_providers {
    aws = "3.44.0"
  }
}

locals {
    timestamp = formatdate("YYYYMMDDhhmmss", timestamp())
    prefix = "route-${local.timestamp}"
}

resource "aws_vpc" "vpc" {
  cidr_block = "10.8.0.0/16"
  tags = {
    Name: "${local.prefix}-default"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "main"
  }
}

resource "aws_default_route_table" "default" {
  default_route_table_id = aws_vpc.vpc.default_route_table_id
  depends_on = [aws_internet_gateway.main]

  route {
    cidr_block = "10.1.1.0/24"
    gateway_id = aws_internet_gateway.main.id
  }

  route {
    cidr_block = "10.1.2.0/24"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table" "r" {
  vpc_id = aws_vpc.vpc.id
  depends_on = [aws_internet_gateway.main]

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "r"
  }
}

resource "aws_route_table" "rr" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "rr"
  }
}

resource "aws_route" "route1" {
  route_table_id = aws_route_table.rr.id
  gateway_id = aws_internet_gateway.main.id
  depends_on = [aws_internet_gateway.main]
  destination_cidr_block = "1.1.1.1/32"
}

resource "aws_route" "route_v6" {
  route_table_id = aws_route_table.rr.id
  gateway_id = aws_internet_gateway.main.id
  depends_on = [aws_internet_gateway.main]
  destination_ipv6_cidr_block = "::/0"
}
