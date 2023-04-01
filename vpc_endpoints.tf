# resource "aws_security_group" "vpce_sg" {
#   name        = "vpce_sg"
#   description = "Allow TLS"
#   vpc_id      = aws_vpc.test-vpc.id

#   ingress {
#     description = "TLS from VPC"
#     from_port   = var.tls_port
#     to_port     = var.tls_port
#     protocol    = var.protocol
#     cidr_blocks = [var.vpc_cidr]
#   }

#   egress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = [var.vpc_cidr]
#   }

#   tags = {
#     Name = "aajolly_vpce_sg"
#   }
# }
# data "aws_subnets" "private" {
#   filter {
#     name   = "vpc-id"
#     values = ["aws_vpc.test-vpc.id"]
#   }
#   filter {
#     name   = "tag:Name"
#     values = ["private*"]
#   }
# }
# resource "aws_vpc_endpoint" "ec2" {
#   for_each          = toset(data.aws_subnets.private.ids)
#   vpc_id            = aws_vpc.test-vpc.id
#   service_name      = "com.amazonaws.ap-southeast-2.ec2"
#   vpc_endpoint_type = "Interface"
#   subnet_ids        = each.value
#   security_group_ids = [
#     aws_security_group.vpce_sg.id,
#   ]

#   private_dns_enabled = true
# }