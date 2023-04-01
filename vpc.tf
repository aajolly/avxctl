# # Create VPC in ap-southeast-2
# locals {
#   project_name = "aajolly"
# }
# resource "aws_vpc" "test-vpc" {
#   provider             = aws.aws-syd
#   cidr_block           = var.vpc_cidr
#   enable_dns_support   = true
#   enable_dns_hostnames = true
#   tags = {
#     Name = "${local.project_name}-test-vpc"
#   }
# }
# #------------------------------------------
# #### Internet Gateway
# #------------------------------------------
# resource "aws_internet_gateway" "igw" {
#   provider = aws.aws-syd
#   vpc_id   = aws_vpc.test-vpc.id
#   tags = {
#     Name = "${local.project_name}-igw"
#   }
# }

# data "aws_availability_zones" "az" {
#   provider = aws.aws-syd
#   state    = "available"
# }
# #------------------------------------------
# #### Subnets
# #------------------------------------------
# resource "aws_subnet" "public_subnet" {
#   provider          = aws.aws-syd
#   count             = 3
#   availability_zone = data.aws_availability_zones.az.names[count.index]
#   vpc_id            = aws_vpc.test-vpc.id
#   cidr_block        = cidrsubnet(aws_vpc.test-vpc.cidr_block, 8, count.index)
#   tags = {
#     Name  = "${local.project_name}-public_subnet_${count.index + 1}"
#     Reach = "Public"
#   }
# }

# resource "aws_subnet" "private_subnet" {
#   provider          = aws.aws-syd
#   count             = 3
#   availability_zone = data.aws_availability_zones.az.names[count.index]
#   vpc_id            = aws_vpc.test-vpc.id
#   cidr_block        = cidrsubnet(aws_vpc.test-vpc.cidr_block, 8, count.index + 4)
#   tags = {
#     Name  = "${local.project_name}-private_subnet_${count.index + 1}"
#     Reach = "Private"
#   }
# }
# #------------------------------------------
# #### Elastic IPs
# #------------------------------------------
# resource "aws_eip" "eip" {
#   count = 3
#   vpc   = true
#   tags = {
#     Name = "${local.project_name}-NAT_GW_EIP_${count.index + 1}"
#   }
#   depends_on = [aws_internet_gateway.igw]
# }
# #------------------------------------------
# #### NAT Gateways
# #------------------------------------------

# resource "aws_nat_gateway" "nat-gw" {
#   count         = 3
#   allocation_id = aws_eip.eip[count.index].id
#   subnet_id     = aws_subnet.public_subnet[count.index].id
#   tags = {
#     Name = "${local.project_name}-NAT_GW_${count.index + 1}"
#   }
# }
# #------------------------------------------
# #### Route Tables
# #------------------------------------------

# resource "aws_route_table" "public_rt" {
#   provider = aws.aws-syd
#   count    = 3
#   vpc_id   = aws_vpc.test-vpc.id
#   tags = {
#     Name = "${local.project_name}-public_rt_${count.index + 1}"
#   }
#   route {
#     cidr_block = var.internet-cidr
#     gateway_id = aws_internet_gateway.igw.id
#   }
# }

# resource "aws_route_table" "private_rt" {
#   provider = aws.aws-syd
#   count    = 3
#   vpc_id   = aws_vpc.test-vpc.id
#   tags = {
#     Name = "${local.project_name}-private_rt_${count.index + 1}"
#   }
#   route {
#     cidr_block     = var.internet-cidr
#     nat_gateway_id = aws_nat_gateway.nat-gw[count.index].id
#   }
# }
# #------------------------------------------
# #### Route Table Associations
# #------------------------------------------
# resource "aws_route_table_association" "public_rt_assoc" {
#   count          = 3
#   subnet_id      = aws_subnet.public_subnet[count.index].id
#   route_table_id = aws_route_table.public_rt[count.index].id
# }

# resource "aws_route_table_association" "private_rt_assoc" {
#   count          = 3
#   subnet_id      = aws_subnet.private_subnet[count.index].id
#   route_table_id = aws_route_table.private_rt[count.index].id
# }