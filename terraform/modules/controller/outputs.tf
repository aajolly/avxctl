output "avx_ctrl_private_ip" {
  value       = aws_instance.aviatrixcontroller.private_ip
  description = "Private IP for Aviatrix Controller"
}

output "avx_ctrl_public_ip" {
  value       = aws_eip.avx_ctrl_eip.public_ip
  description = "List of public IPs for aviatrix conroller"
}

output "vpc_id" {
  value       = aws_vpc.avx-mgmt-vpc.id
  description = "VPC where Aviatrix Controller was built"
}

output "subnet_id" {
  value       = aws_subnet.public_subnet[0].id
  description = "Subnet where Aviatrix Controller was built"
}

output "security_group_id" {
  value       = aws_security_group.avx-sg.id
  description = "Security group id used by Aviatrix Controller"
}