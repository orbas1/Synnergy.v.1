terraform {
  required_version = ">= 1.3"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

locals {
  common_tags = {
    Project     = "synnergy-network"
    Environment = "stage77"
  }
}

variable "region" {
  description = "AWS region to deploy Synnergy Network infrastructure."
  type        = string
  default     = "us-east-1"
}

# Base networking
resource "aws_vpc" "synnergy" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "synnergy-vpc"
  }
}

data "aws_availability_zones" "available" {}

resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.synnergy.id
  cidr_block              = cidrsubnet(aws_vpc.synnergy.cidr_block, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = {
    Name = "synnergy-public-${count.index}"
  }
}

resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.synnergy.id
  cidr_block        = cidrsubnet(aws_vpc.synnergy.cidr_block, 8, count.index + 10)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = "synnergy-private-${count.index}"
  }
}

# Security group for node communication
resource "aws_security_group" "node" {
  name        = "synnergy-node-sg"
  description = "Security group for Synnergy nodes"
  vpc_id      = aws_vpc.synnergy.id

  ingress {
    description = "P2P communication"
    from_port   = 30303
    to_port     = 30303
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "RPC"
    from_port   = 8080
    to_port     = 8080
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

resource "aws_security_group" "wallet" {
  name        = "synnergy-wallet-sg"
  description = "Security group for wallet servers"
  vpc_id      = aws_vpc.synnergy.id

  ingress {
    description = "Wallet HTTPS"
    from_port   = 8081
    to_port     = 8081
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

resource "aws_kms_key" "ledger" {
  description             = "Synnergy ledger encryption"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  tags                    = local.common_tags
}

resource "aws_cloudwatch_log_group" "synnergy" {
  name              = "/synnergy/network"
  retention_in_days = 30
  kms_key_id        = aws_kms_key.ledger.arn
  tags              = local.common_tags
}

# Launch template for nodes. User data will bootstrap ansible.
resource "aws_launch_template" "node" {
  name_prefix   = "synnergy-node-"
  image_id      = var.ami_id
  instance_type = "t3.medium"

  user_data = base64encode(<<-EOT
    #!/bin/bash
    apt-get update && apt-get install -y python3
    echo "$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)" >> /tmp/node_ip
    /usr/local/bin/synnergy bootstrap --gas-report --telemetry-endpoint ${var.telemetry_endpoint}
  EOT
  )

  iam_instance_profile {
    name = aws_iam_instance_profile.node.name
  }

  vpc_security_group_ids = [aws_security_group.node.id]
  tag_specifications {
    resource_type = "instance"
    tags          = local.common_tags
  }
}

# Auto scaling group for fault-tolerant nodes
resource "aws_autoscaling_group" "nodes" {
  name                = "synnergy-nodes"
  max_size            = 3
  min_size            = 1
  desired_capacity    = 1
  vpc_zone_identifier = aws_subnet.public[*].id
  launch_template {
    id      = aws_launch_template.node.id
    version = "$Latest"
  }
  target_group_arns = [aws_lb_target_group.synnergy.arn]
  tag {
    key                 = "Name"
    value               = "synnergy-node"
    propagate_at_launch = true
  }
}

variable "ami_id" {
  description = "AMI ID for node instances"
  type        = string
}

variable "telemetry_endpoint" {
  description = "Endpoint for forwarding OpenTelemetry traces"
  type        = string
  default     = "https://telemetry.synnergy.net"
}

variable "db_username" {
  description = "Master username for the Synnergy wallet database"
  type        = string
}

variable "db_password" {
  description = "Master password for the Synnergy wallet database"
  type        = string
  sensitive   = true
}

resource "aws_db_subnet_group" "wallet" {
  name       = "synnergy-wallet-subnets"
  subnet_ids = aws_subnet.private[*].id
  tags       = local.common_tags
}

resource "aws_rds_cluster" "wallet" {
  cluster_identifier      = "synnergy-wallet"
  engine                  = "aurora-postgresql"
  engine_version          = "15.4"
  master_username         = var.db_username
  master_password         = var.db_password
  db_subnet_group_name    = aws_db_subnet_group.wallet.name
  storage_encrypted       = true
  kms_key_id              = aws_kms_key.ledger.arn
  backup_retention_period = 7
  preferred_backup_window = "03:00-04:00"
  skip_final_snapshot     = false
  deletion_protection     = true
  vpc_security_group_ids  = [aws_security_group.wallet.id]
  tags                    = local.common_tags
}

resource "aws_lb" "synnergy" {
  name               = "synnergy-alb"
  load_balancer_type = "application"
  subnets            = aws_subnet.public[*].id
  security_groups    = [aws_security_group.node.id]
  idle_timeout       = 60
  enable_deletion_protection = true
  tags = local.common_tags
}

resource "aws_lb_target_group" "synnergy" {
  name        = "synnergy-node-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.synnergy.id
  health_check {
    path                = "/healthz"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200"
  }
  tags = local.common_tags
}

resource "aws_lb_listener" "synnergy" {
  load_balancer_arn = aws_lb.synnergy.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = var.listener_certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.synnergy.arn
  }
}

variable "listener_certificate_arn" {
  description = "ACM certificate ARN for the public listener"
  type        = string
}

resource "aws_iam_role" "node" {
  name = "synnergy-node"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
  tags = local.common_tags
}

resource "aws_iam_role_policy" "node" {
  name = "synnergy-node"
  role = aws_iam_role.node.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["logs:CreateLogStream", "logs:PutLogEvents"]
        Resource = "${aws_cloudwatch_log_group.synnergy.arn}:*"
      },
      {
        Effect   = "Allow"
        Action   = ["kms:Encrypt", "kms:Decrypt", "kms:GenerateDataKey"]
        Resource = aws_kms_key.ledger.arn
      }
    ]
  })
}

resource "aws_iam_instance_profile" "node" {
  name = "synnergy-node"
  role = aws_iam_role.node.name
}

output "asg_name" {
  description = "Name of the autoscaling group hosting Synnergy nodes."
  value       = aws_autoscaling_group.nodes.name
}

output "load_balancer_dns" {
  description = "DNS name of the Synnergy application load balancer"
  value       = aws_lb.synnergy.dns_name
}

output "wallet_cluster_endpoint" {
  description = "Reader endpoint for the wallet Aurora cluster"
  value       = aws_rds_cluster.wallet.reader_endpoint
}
