#--------------------------------------------------
## IAM Resources - Policies, Roles
#--------------------------------------------------

resource "aws_iam_policy" "avx_ctrl_role_ec2_policy" {
  name = "aviatrix-role-ec2-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["sts:AssumeRole"]
        Effect   = "Allow"
        Resource = "arn:aws::iam::*:role/aviatrix-*"
      },
      {
        Action = [
          "aws-marketplace:MeterUsage",
          "s3:GetBucketLocation"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_policy" "avx_ctrl_role_app_policy" {
  name = "${local.name_prefix}-role-app-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ec2:Describe*",
          "ec2:Get*",
          "ec2:Search*",
          "elasticloadbalancing:Describe*",
          "route53:List*",
          "route53:Get*",
          "sqs:Get*",
          "sqs:List*",
          "sns:List*",
          "s3:List*",
          "s3:Get*",
          "iam:List*",
          "iam:Get*",
          "directconnect:Describe*",
          "guardduty:Get*",
          "guardduty:List*",
          "ram:Get*",
          "ram:List*",
          "networkmanager:Get*",
          "networkmanager:List*"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "ec2:CreateVpc",
          "ec2:DeleteVpc",
          "ec2:ModifyVpcAttribute",
          "ec2:CreateNetworkAclEntry",
          "ec2:ReplaceNetworkAclEntry",
          "ec2:DeleteNetworkAclEntry",
          "ec2:AssociateVpcCidrBlock",
          "ec2:AssociateSubnetCidrBlock",
          "ec2:CreateSubnet",
          "ec2:DeleteSubnet",
          "ec2:ModifySubnetAttribute",
          "ec2:*InternetGateway*",
          "ec2:*Route*",
          "ec2:*Instance*",
          "ec2:*SecurityGroup*",
          "ec2:*Address*",
          "ec2:*NetworkInterface*",
          "ec2:CreateKeyPair",
          "ec2:DeleteKeyPair",
          "ec2:CreateTags",
          "ec2:DeleteTags"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "ec2:CreateCustomerGateway",
          "ec2:DeleteCustomerGateway",
          "ec2:CreateVpnConnection",
          "ec2:DeleteVpnConnection",
          "ec2:CreateVpcPeeringConnection",
          "ec2:AcceptVpcPeeringConnection",
          "ec2:DeleteVpcPeeringConnection",
          "ec2:EnableVgwRoutePropagation",
          "ec2:DisableVgwRoutePropagation"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "ec2:AssociateTransitGatewayRouteTable",
          "ec2:AcceptTransitGatewayVpcAttachment",
          "ec2:CreateTransitGateway",
          "ec2:CreateTransitGatewayRoute",
          "ec2:CreateTransitGatewayRouteTable",
          "ec2:CreateTransitGatewayVpcAttachment",
          "ec2:DeleteTransitGateway",
          "ec2:DeleteTransitGatewayRoute",
          "ec2:DeleteTransitGatewayRouteTable",
          "ec2:DeleteTransitGatewayVpcAttachment",
          "ec2:DisableTransitGatewayRouteTablePropagation",
          "ec2:DisassociateTransitGatewayRouteTable",
          "ec2:EnableTransitGatewayRouteTablePropagation",
          "ec2:ExportTransitGatewayRoutes",
          "ec2:ModifyTransitGatewayVpcAttachment",
          "ec2:RejectTransitGatewayVpcAttachment",
          "ec2:ReplaceTransitGatewayRoute",
          "ec2:EnableRoutePropagation"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "ram:CreateResourceShare",
          "ram:DeleteResourceShare",
          "ram:UpdateResourceShare",
          "ram:AssociateResourceShare",
          "ram:DisassociateResourceShare",
          "ram:TagResource",
          "ram:UntagResource",
          "ram:AcceptResourceShareInvitation",
          "ram:EnableSharingWithAwsOrganization"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "directconnect:CreateDirectConnectGateway",
          "directconnect:CreateDirectConnectGatewayAssociation",
          "directconnect:CreateDirectConnectGatewayAssociationProposal",
          "directconnect:DeleteDirectConnectGateway",
          "directconnect:DeleteDirectConnectGatewayAssociation",
          "directconnect:DeleteDirectConnectGatewayAssociationProposal",
          "directconnect:AcceptDirectConnectGatewayAssociationProposal"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "sqs:AddPermission",
          "sqs:ChangeMessageVisibility",
          "sqs:CreateQueue",
          "sqs:DeleteMessage",
          "sqs:DeleteQueue",
          "sqs:PurgeQueue",
          "sqs:ReceiveMessage",
          "sqs:RemovePermission",
          "sqs:SendMessage",
          "sqs:SetQueueAttributes",
          "sqs:TagQueue"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "s3:CreateBucket",
          "s3:DeleteBucket",
          "s3:ListBucket",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "iam:PassRole",
          "iam:AddRoleToInstanceProfile",
          "iam:CreateInstanceProfile",
          "iam:DeleteInstanceProfile",
          "iam:RemoveRoleFromInstanceProfile",
          "iam:CreateServiceLinkedRole"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "iam:DeletePolicyVersion",
          "iam:CreatePolicyVersion"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "elasticloadbalancing:*",
          "route53:ChangeResourceRecordSets",
          "ec2:*Volume*",
          "ec2:*Snapshot*",
          "ec2:*TransitGatewayPeeringAttachment",
          "guardduty:*",
          "globalaccelerator:*",
          "networkmanager:*"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role" "avx_ctrl_role_ec2" {
  name                = "${local.name_prefix}-role-ec2"
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
    Name = "${local.name_prefix}-role-ec2"
  }
}

resource "aws_iam_role" "avx_ctrl_role_app" {
  name        = "${local.name_prefix}-role-app"
  path        = "/"
  description = "Aviatrix App Role"
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  ]
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
    Name = "${local.name_prefix}-role-app"
  }
}

resource "aws_iam_instance_profile" "avx_ctrl_ec2_profile" {
  name = "${local.name_prefix}-ec2-profile"
  role = aws_iam_role.avx_ctrl_role_ec2.name
}

#-------------------------------------------------------------
## VPC Resources - VPC, Subnets, Route Tables, Routes etc
#-------------------------------------------------------------
resource "aws_vpc" "avx-mgmt-vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "${local.name_prefix}-mgmt-vpc"
  }
}

#------------------------------------------
#### Internet Gateway
#------------------------------------------
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.avx-mgmt-vpc.id
  tags = {
    Name = "${local.name_prefix}-igw"
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
  vpc_id            = aws_vpc.avx-mgmt-vpc.id
  cidr_block        = cidrsubnet(aws_vpc.avx-mgmt-vpc.cidr_block, 4, count.index)
  tags = {
    Name  = "${local.name_prefix}-public_subnet_${count.index + 1}"
    Reach = "Public"
  }
}
#------------------------------------------
#### Elastic IPs
#------------------------------------------
resource "aws_eip" "avx_ctrl_eip" {
  vpc = true
  tags = {
    Name = "${local.name_prefix}-ctrl-eip"
  }
  depends_on = [aws_internet_gateway.igw]
}

resource "aws_eip_association" "avx_ctrl_eip_assoc" {
  instance_id   = aws_instance.aviatrixcontroller.id
  allocation_id = aws_eip.avx_ctrl_eip.id
}

resource "aws_network_interface" "avx_ctrl_eni" {
  subnet_id       = aws_subnet.public_subnet[0].id
  security_groups = [aws_security_group.avx-sg.id]
  tags = {
    Name      = "${local.name_prefix}-controller-interface"
    Createdby = "Terraform+Aviatrix"
  }
}
#------------------------------------------
#### Route Tables
#------------------------------------------

resource "aws_route_table" "public_rt" {
  count  = 3
  vpc_id = aws_vpc.avx-mgmt-vpc.id
  tags = {
    Name = "${local.name_prefix}-public_rt_${count.index + 1}"
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
resource "aws_security_group" "avx-sg" {
  description = "Aviatrix - Allow HTTPS to Controller"
  vpc_id      = aws_vpc.avx-mgmt-vpc.id

  ingress {
    description = "Allow Access from My Public IP"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [local.my_public_ip]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [var.internet_cidr]
  }

  tags = {
    Name = "${local.name_prefix}SecurityGroup"
  }
}
#------------------------------------------
#### Aviatrix Controller
#------------------------------------------
resource "aws_instance" "aviatrixcontroller" {
  ami                     = local.ami_id
  instance_type           = var.instance_type
  key_name                = "aajolly-apse2"
  iam_instance_profile    = aws_iam_role.avx_ctrl_role_ec2.name
  disable_api_termination = true

  network_interface {
    network_interface_id = aws_network_interface.avx_ctrl_eni.id
    device_index         = 0
  }

  root_block_device {
    volume_size = var.root_volume_size
    volume_type = var.root_volume_type
    encrypted   = var.ebs-encryption
  }

  tags = {
    Name      = "${local.name_prefix}-controller"
    Createdby = "Terraform+Aviatrix"
  }

  lifecycle {
    ignore_changes = [
      ami
    ]
  }
  user_data = <<EOF
#!/bin/bash -ex
sudo service ssh stop
EOF
}

resource "null_resource" "this" {
  provisioner "local-exec" {
    command     = <<-EOT
    #expects json file, also private_ip of ctrl, pub_ip of copilot
    # waiting for the 200 response
        while [[ "$(curl -s -o /dev/null -w '%%{http_code}' https://${aws_instance.aviatrixcontroller.public_ip} --insecure)" != "200" ]]; do sleep 10; echo "waiting for controller"; done
    # in the beginning autheticate with username and password
        while [[ "$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=login&username=admin&password=${aws_instance.aviatrixcontroller.private_ip}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure | jq -r .return)" != "true" ]]; do sleep 10; echo "Waiting for CID   "; done
        sleep 60
        init_ctrl_auth=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=login&username=admin&password=${aws_instance.aviatrixcontroller.private_ip}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
    # init_ctrl_auth=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=login&username=admin&password=${var.admin_password}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        CID=$(echo $init_ctrl_auth | jq -r .CID)
        echo $CID
    #auth with CID , set password, set recovery email, defaultemail : ace.lab@aviatrix.com tfvars may already exist, can reuse
        set_password=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=edit_account_user&CID=$CID&account_name=admin&username=admin&password=${var.admin_password}&what=password&email=${var.admin_email}&old_password=${aws_instance.aviatrixcontroller.private_ip}&new_password=${var.admin_password}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        echo $set_password
        add_email=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=add_admin_email_addr&CID=$CID&admin_email=${var.admin_email}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        echo $add_email
    #should be over writerable, curr set to pod name (pod130), just dont set controller label if not passed in
        set_ctrl_label=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=set_controller_name&CID=$CID&controller_name=${var.ctrl_name}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        echo $set_ctrl_label
    #if not passed in, don't do it
        set_ctrl_license=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=setup_customer_id&CID=$CID&customer_id=${var.ctrl_customer_id}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        echo $set_ctrl_license
    #create that service acc in ctrl
        add_copilot_service_account=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d "action=add_account_user&CID=$CID&account_name=copilot&username=copilot&password=${var.admin_password}&email=${var.admin_email}&groups=admin" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        echo $add_copilot_service_account
    #upgrade
        upgrade_ctrl=$(curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' -m 360 -d "action=upgrade&CID=$CID&version=${var.ctrl_version}" https://${aws_instance.aviatrixcontroller.public_ip}/v1/api --insecure)
        echo $upgrade_ctrl
    #don't sleep, know when it's done
        sleep 150
    EOT
    interpreter = ["/bin/bash", "-c"]
  }
}
