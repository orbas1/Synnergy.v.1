terraform {
  required_version = ">= 1.5"
  backend "s3" {
    bucket         = var.state_bucket
    key            = "synnergy/network/terraform.tfstate"
    region         = var.region
    dynamodb_table = var.state_lock_table
    encrypt        = true
  }

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

variable "region" {
  description = "AWS region to deploy Synnergy Network infrastructure."
  type        = string
  default     = "us-east-1"
}

variable "state_bucket" {
  description = "Name of the S3 bucket used for the Terraform remote state."
  type        = string
}

variable "state_lock_table" {
  description = "Name of the DynamoDB table used for Terraform state locking."
  type        = string
}

variable "project" {
  description = "Project tag applied to all infrastructure resources."
  type        = string
  default     = "synnergy"
}

variable "environment" {
  description = "Deployment environment tag (e.g. dev, staging, prod)."
  type        = string
  default     = "prod"
}

variable "node_instance_type" {
  description = "EC2 instance type for Synnergy nodes."
  type        = string
  default     = "m6i.large"
}

variable "wallet_instance_type" {
  description = "EC2 instance type for wallet API nodes."
  type        = string
  default     = "t3.medium"
}

variable "db_password" {
  description = "Password for the wallet audit database."
  type        = string
  sensitive   = true
}

locals {
  common_tags = {
    Project     = var.project
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

resource "aws_kms_key" "config" {
  description             = "Synnergy runtime configuration encryption"
  deletion_window_in_days = 7
  tags                    = local.common_tags
}

resource "aws_cloudwatch_log_group" "synnergy" {
  name              = "/aws/synnergy/network"
  retention_in_days = 30
  kms_key_id        = aws_kms_key.config.arn
  tags              = local.common_tags
}

resource "aws_vpc" "synnergy" {
  cidr_block = "10.40.0.0/16"
  tags       = merge(local.common_tags, { Name = "synnergy-vpc" })
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.synnergy.id
  tags   = merge(local.common_tags, { Name = "synnergy-igw" })
}

data "aws_availability_zones" "available" {}

resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.synnergy.id
  cidr_block              = cidrsubnet(aws_vpc.synnergy.cidr_block, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = merge(local.common_tags, {
    Name = "synnergy-public-${count.index}"
    Tier = "public"
  })
}

resource "aws_subnet" "private" {
  count                   = 2
  vpc_id                  = aws_vpc.synnergy.id
  cidr_block              = cidrsubnet(aws_vpc.synnergy.cidr_block, 8, count.index + 10)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = false
  tags = merge(local.common_tags, {
    Name = "synnergy-private-${count.index}"
    Tier = "private"
  })
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.synnergy.id
  tags   = merge(local.common_tags, { Name = "synnergy-public" })
}

resource "aws_route" "internet" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw.id
}

resource "aws_route_table_association" "public" {
  count          = length(aws_subnet.public)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "node" {
  name        = "synnergy-node"
  description = "Security group for Synnergy blockchain nodes"
  vpc_id      = aws_vpc.synnergy.id

  ingress {
    description = "CLI and RPC access"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "P2P networking"
    from_port   = 30303
    to_port     = 30303
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, { Name = "synnergy-node" })
}

resource "aws_security_group" "wallet" {
  name        = "synnergy-wallet"
  description = "Security group for wallet API nodes"
  vpc_id      = aws_vpc.synnergy.id

  ingress {
    description = "Wallet HTTPS"
    from_port   = 8090
    to_port     = 8090
    protocol    = "tcp"
    security_groups = [aws_security_group.node.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, { Name = "synnergy-wallet" })
}

resource "aws_lb" "node" {
  name               = "synnergy-node"
  load_balancer_type = "network"
  subnets            = aws_subnet.public[*].id
  tags               = local.common_tags
}

resource "aws_lb_target_group" "node_rpc" {
  name     = "syn-node-rpc"
  port     = 8080
  protocol = "TCP"
  vpc_id   = aws_vpc.synnergy.id
  health_check {
    port                = "8080"
    protocol            = "TCP"
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
  tags = local.common_tags
}

resource "aws_lb_listener" "rpc" {
  load_balancer_arn = aws_lb.node.arn
  port              = 8080
  protocol          = "TCP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.node_rpc.arn
  }
}

resource "aws_launch_template" "node" {
  name_prefix   = "syn-node-"
  image_id      = var.node_ami_id
  instance_type = var.node_instance_type
  key_name      = var.ssh_key_name

  iam_instance_profile {
    name = aws_iam_instance_profile.synnergy.name
  }

  network_interfaces {
    security_groups = [aws_security_group.node.id]
  }

  user_data = base64encode(<<-EOT
    #!/bin/bash
    set -euo pipefail
    amazon-linux-extras install epel -y
    yum install -y jq
    aws ssm get-parameter --with-decryption --name "/synnergy/node/config" \
      --query Parameter.Value --output text > /etc/synnergy/synnergy.yaml
    systemctl enable synnergy
    systemctl start synnergy
  EOT
  )

  tag_specifications {
    resource_type = "instance"
    tags          = merge(local.common_tags, { Role = "node" })
  }
}

variable "node_ami_id" {
  description = "AMI ID for Synnergy node instances"
  type        = string
}

variable "ssh_key_name" {
  description = "SSH key name to attach to launched instances."
  type        = string
}

resource "aws_autoscaling_group" "node" {
  name                = "synnergy-node"
  max_size            = 5
  min_size            = 3
  desired_capacity    = 3
  vpc_zone_identifier = aws_subnet.public[*].id
  launch_template {
    id      = aws_launch_template.node.id
    version = "$Latest"
  }
  target_group_arns = [aws_lb_target_group.node_rpc.arn]
  metrics_granularity = "1Minute"
  enabled_metrics      = ["GroupDesiredCapacity", "GroupInServiceInstances"]
  tag {
    key                 = "Name"
    value               = "synnergy-node"
    propagate_at_launch = true
  }
}

resource "aws_launch_template" "wallet" {
  name_prefix   = "syn-wallet-"
  image_id      = var.wallet_ami_id
  instance_type = var.wallet_instance_type

  iam_instance_profile {
    name = aws_iam_instance_profile.synnergy.name
  }

  network_interfaces {
    associate_public_ip_address = false
    subnet_id                   = aws_subnet.private[0].id
    security_groups             = [aws_security_group.wallet.id]
  }

  user_data = base64encode(<<-EOT
    #!/bin/bash
    set -euo pipefail
    aws secretsmanager get-secret-value --secret-id synnergy-wallet \
      --query SecretString --output text > /etc/synnergy/wallet.json
    systemctl enable synnergy-wallet
    systemctl start synnergy-wallet
  EOT
  )

  tag_specifications {
    resource_type = "instance"
    tags          = merge(local.common_tags, { Role = "wallet" })
  }
}

variable "wallet_ami_id" {
  description = "AMI ID for wallet instances"
  type        = string
}

resource "aws_autoscaling_group" "wallet" {
  name                = "synnergy-wallet"
  max_size            = 4
  min_size            = 2
  desired_capacity    = 2
  vpc_zone_identifier = aws_subnet.private[*].id
  launch_template {
    id      = aws_launch_template.wallet.id
    version = "$Latest"
  }
  tag {
    key                 = "Name"
    value               = "synnergy-wallet"
    propagate_at_launch = true
  }
}

resource "aws_db_subnet_group" "wallet" {
  name       = "synnergy-wallet"
  subnet_ids = aws_subnet.private[*].id
  tags       = local.common_tags
}

resource "aws_db_instance" "wallet_audit" {
  identifier              = "synnergy-wallet"
  allocated_storage       = 20
  engine                  = "postgres"
  engine_version          = "16.3"
  instance_class          = "db.t3.micro"
  db_name                 = "synnergy"
  username                = "synnergy"
  password                = var.db_password
  skip_final_snapshot     = true
  storage_encrypted       = true
  kms_key_id              = aws_kms_key.config.arn
  backup_retention_period = 7
  db_subnet_group_name    = aws_db_subnet_group.wallet.name
  vpc_security_group_ids  = [aws_security_group.wallet.id]
  tags                    = local.common_tags
}

resource "aws_ssm_parameter" "node_config" {
  name        = "/synnergy/node/config"
  type        = "SecureString"
  key_id      = aws_kms_key.config.arn
  value       = jsonencode({
    logLevel    = "info"
    environment = var.environment
    rpc         = "0.0.0.0:8080"
  })
  overwrite   = true
  description = "Synnergy node runtime configuration"
  tags        = local.common_tags
}

resource "aws_iam_role" "synnergy" {
  name               = "synnergy-runtime"
  assume_role_policy = data.aws_iam_policy_document.ec2.json
  tags               = local.common_tags
}

data "aws_iam_policy_document" "ec2" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_instance_profile" "synnergy" {
  name = "synnergy-runtime"
  role = aws_iam_role.synnergy.name
}

resource "aws_iam_role_policy" "synnergy" {
  name = "synnergy-runtime"
  role = aws_iam_role.synnergy.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["ssm:GetParameter"],
        Effect   = "Allow",
        Resource = aws_ssm_parameter.node_config.arn
      },
      {
        Action   = ["logs:CreateLogStream", "logs:PutLogEvents"],
        Effect   = "Allow",
        Resource = "${aws_cloudwatch_log_group.synnergy.arn}:*"
      },
      {
        Action   = ["secretsmanager:GetSecretValue"],
        Effect   = "Allow",
        Resource = "*"
      }
    ]
  })
}

output "node_load_balancer" {
  description = "DNS name for the Synnergy node network load balancer"
  value       = aws_lb.node.dns_name
}

output "wallet_database_endpoint" {
  description = "Endpoint for the wallet audit database"
  value       = aws_db_instance.wallet_audit.endpoint
}
