terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}
provider "aws" {
  region = var.region
}
#--------------------------------------------------
## IAM Resources - Policies, Roles
#--------------------------------------------------

resource "aws_iam_policy" "avx_ctrl_role_ec2_policy" {
  name = "${local.tool_prefix}-role-ec2-policy-${local.region}"

  policy = data.aws_iam_policy_document.avx_role_ec2_policy.json
}

resource "aws_iam_policy" "avx_ctrl_role_app_policy" {
  name = "${local.tool_prefix}-role-app-policy-${local.region}"

  policy = data.http.avx_role_app_policy.response_body
}

resource "aws_iam_role" "avx_ctrl_role_ec2" {
  name                = "${local.tool_prefix}-role-ec2-${local.region}"
  path                = "/"
  description         = "Aviatrix EC2 Role"
  managed_policy_arns = [aws_iam_policy.avx_ctrl_role_ec2_policy.arn]
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })

  tags = {
    Name = "${local.tool_prefix}-role-ec2-${local.region}"
    Createdby = local.tool_prefix
  }
}

resource "aws_iam_role" "avx_ctrl_role_app" {
  name        = "${local.tool_prefix}-role-app-${local.region}"
  path        = "/"
  description = "Aviatrix App Role"
  managed_policy_arns = [aws_iam_policy.avx_ctrl_role_app_policy.arn]
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          AWS = "arn:aws:iam::${local.account_id}:root"
        }
      },
    ]
  })

  tags = {
    Name = "${local.tool_prefix}-role-app-${local.region}"
    Createdby = local.tool_prefix
  }
}

resource "aws_iam_instance_profile" "avx_ctrl_ec2_profile" {
  name = "${local.tool_prefix}-ec2-profile-${local.region}"
  role = aws_iam_role.avx_ctrl_role_ec2.name
}

#-------------------------------------------------------------
## VPC Resources - VPC, Subnets, Route Tables, Routes etc
#-------------------------------------------------------------
resource "aws_vpc" "avx_mgmt_vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "${local.tool_prefix}-mgmt-vpc"
    Createdby = local.tool_prefix
  }
}

#------------------------------------------
#### Internet Gateway
#------------------------------------------
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.avx_mgmt_vpc.id
  tags = {
    Name = "${local.tool_prefix}-igw"
    Createdby = local.tool_prefix
  }
}

data "aws_availability_zones" "az" {
  state = "available"
}
#------------------------------------------
#### Subnets
#------------------------------------------
resource "aws_subnet" "public_subnet" {
  count             = 3
  availability_zone = data.aws_availability_zones.az.names[count.index]
  vpc_id            = aws_vpc.avx_mgmt_vpc.id
  cidr_block        = cidrsubnet(aws_vpc.avx_mgmt_vpc.cidr_block, 4, count.index)
  tags = {
    Name  = "${local.tool_prefix}-public_subnet_${count.index + 1}"
    Reach = "Public"
    Createdby = local.tool_prefix
  }
}
#------------------------------------------
#### Elastic IPs
#------------------------------------------
resource "aws_eip" "avx_ctrl_eip" {
  vpc = true
  tags = {
    Name = "${local.tool_prefix}-ctrl-eip"
    Createdby = local.tool_prefix
  }
  depends_on = [aws_internet_gateway.igw]
}

resource "aws_eip_association" "avx_ctrl_eip_assoc" {
  instance_id   = aws_instance.aviatrixcontroller.id
  allocation_id = aws_eip.avx_ctrl_eip.id
}

resource "aws_network_interface" "avx_ctrl_eni" {
  subnet_id       = aws_subnet.public_subnet[0].id
  security_groups = [aws_security_group.avx_sg.id]
  tags = {
    Name      = "${local.tool_prefix}-controller-interface"
    Createdby = local.tool_prefix
  }
}
#------------------------------------------
#### Route Tables
#------------------------------------------

resource "aws_route_table" "public_rt" {
  count  = 3
  vpc_id = aws_vpc.avx_mgmt_vpc.id
  tags = {
    Name = "${local.tool_prefix}-public_rt_${count.index + 1}"
    Createdby = local.tool_prefix
  }
  route {
    cidr_block = var.internet_cidr
    gateway_id = aws_internet_gateway.igw.id
  }
}
#------------------------------------------
#### Route Table Associations
#------------------------------------------
resource "aws_route_table_association" "public_rt_assoc" {
  count          = 3
  subnet_id      = aws_subnet.public_subnet[count.index].id
  route_table_id = aws_route_table.public_rt[count.index].id
}
#------------------------------------------
#### Security Group
#------------------------------------------
resource "aws_security_group" "avx_sg" {
  description = "Aviatrix - Allow HTTPS to Controller"
  vpc_id      = aws_vpc.avx_mgmt_vpc.id

  tags = {
    Name = "${local.tool_prefix}SecurityGroup"
    Createdby = local.tool_prefix
  }
}
resource "aws_vpc_security_group_ingress_rule" "avx_ctrl_ingress1" {
  description = "Allow Ingress from local machine"
  security_group_id = aws_security_group.avx_sg.id
  cidr_ipv4 = local.my_public_ip
  from_port = 443
  to_port = 443
  ip_protocol = "tcp"
}

resource "aws_vpc_security_group_egress_rule" "avx_ctrl_egress1" {
  description = "Allow Egress from Aviatrix Controller"
  security_group_id = aws_security_group.avx_sg.id
  cidr_ipv4 = var.internet_cidr
  ip_protocol = "-1"
}
#------------------------------------------
#### Aviatrix Controller
#------------------------------------------
resource "aws_instance" "aviatrixcontroller" {
  ami                     = local.ami_id
  instance_type           = var.instance_type
  key_name                = var.keypair
  iam_instance_profile    = aws_iam_instance_profile.avx_ctrl_ec2_profile.name
  disable_api_termination = true

  network_interface {
    network_interface_id = aws_network_interface.avx_ctrl_eni.id
    device_index         = 0
  }

  root_block_device {
    volume_size = var.root_volume_size
    volume_type = var.root_volume_type
    encrypted   = var.ebs_encryption
  }

  tags = {
    Name      = "${local.tool_prefix}-ctrl-${var.ctrl_version}"
    Createdby = local.tool_prefix
  }

  lifecycle {
    ignore_changes = [
      ami
    ]
  }
}

#------------------------------------------
#### Aviatrix CoPilot
#------------------------------------------
resource "aws_ebs_volume" "copilot" {
  count = var.ctrl_version >= 6.8 ? 0 : 1
  availability_zone = "${data.aws_region.current.name}a"
  encrypted         = true
  type              = "gp2"
  size              = 30
}

module "copilot_build_aws" {
  count = var.ctrl_version >= 6.8 ? 0 : 1
  source                = "github.com/AviatrixSystems/terraform-modules-copilot.git//copilot_build_aws"
  copilot_name          = "${local.tool_prefix}-copilot"
  use_existing_keypair  = true
  keypair               = var.keypair
  controller_public_ip  = aws_eip.avx_ctrl_eip.public_ip
  controller_private_ip = aws_instance.aviatrixcontroller.private_ip
  instance_type         = "t3.xlarge"
  use_existing_vpc      = true
  vpc_id                = aws_vpc.avx_mgmt_vpc.id
  subnet_id             = aws_subnet.public_subnet[0].id

  allowed_cidrs = {
    "tcp_cidrs" = {
      protocol = "tcp"
      port     = "443"
      cidrs    = ["0.0.0.0/0"]
    }
    "udp_cidrs_1" = {
      protocol = "udp"
      port     = "5000"
      cidrs    = ["0.0.0.0/0"]
    }
    "udp_cidrs_2" = {
      protocol = "udp"
      port     = "31283"
      cidrs    = ["0.0.0.0/0"]
    }
  }
  additional_volumes = {
    "one" = {
      device_name = "/dev/sda2"
      volume_id   = aws_ebs_volume.copilot[count.index].id
    }
  }
}