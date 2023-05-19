output "avx_ctrl_private_ip" {
  value       = aws_instance.aviatrixcontroller.private_ip
  description = "Private IP for Aviatrix Controller"
}

output "avx_ctrl_public_ip" {
  value       = aws_eip.avx_ctrl_eip.public_ip
  description = "List of public IPs for aviatrix conroller"
}

output "avxctl_ctrl_vpc_id" {
  value       = aws_vpc.avx_mgmt_vpc.id
  description = "VPC where Aviatrix Controller was built"
}

output "avxctl_ctrl_subnet_id" {
  value       = aws_subnet.public_subnet[0].id
  description = "Subnet where Aviatrix Controller was built"
}

output "avxctl_ctrl_subnet_cidr" {
  value       = aws_subnet.public_subnet[0].cidr_block
  description = "Subnet CIDR where Aviatrix Controller was built"
}

output "security_group_id" {
  value       = aws_security_group.avx_sg.id
  description = "Security group id used by Aviatrix Controller"
}
# output "copilot_public_ip" {
#   value = module.copilot_build_aws[count.index].public_ip
# }

output "aws_role_arn" {
  value = aws_iam_role.avx_ctrl_role_app.arn
}

output "aws_role_ec2" {
  value = aws_iam_role.avx_ctrl_role_ec2.arn
}