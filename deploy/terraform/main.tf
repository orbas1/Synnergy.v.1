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

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
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
  EOT
  )
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

output "asg_name" {
  description = "Name of the autoscaling group hosting Synnergy nodes."
  value       = aws_autoscaling_group.nodes.name
}
