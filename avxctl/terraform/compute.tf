# data "aws_ami" "al2" {
#   most_recent = true

#   filter {
#     name   = "owner-alias"
#     values = ["amazon"]
#   }

#   filter {
#     name   = "name"
#     values = ["amzn-ami-hvm*"]
#   }

#   filter {
#     name   = "root-device-type"
#     values = ["ebs"]
#   }

#   filter {
#     name   = "virtualization-type"
#     values = ["hvm"]
#   }
# }

# resource "aws_security_group" "sg1" {
#   name        = "aajolly_ec2_sg"
#   description = "Allow TLS"
#   vpc_id      = aws_vpc.test-vpc.id

#   ingress {
#     description = "TLS from Internet"
#     from_port   = var.tls_port
#     to_port     = var.tls_port
#     protocol    = var.protocol
#     cidr_blocks = [var.internet-cidr]
#   }

#   egress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = [var.internet-cidr]
#   }

#   tags = {
#     Name = "aajolly_ec2_sg"
#   }
# }
# resource "aws_instance" "aajolly-test-ec2" {
#   ami                  = data.aws_ami.al2.id
#   instance_type        = var.instance_type
#   iam_instance_profile = aws_iam_instance_profile.aajolly_test_profile.name
#   root_block_device {
#     volume_size = "20"
#     volume_type = "gp2"
#     encrypted   = true
#     tags = {
#       Name = "aajolly-test-ec2"
#     }
#   }
#   vpc_security_group_ids = [aws_security_group.sg1.id]
#   subnet_id              = aws_subnet.private_subnet[0].id

#   tags = {
#     Name = "aajolly-test-ec2"
#   }
# }