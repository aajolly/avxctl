output "vpc_id" {
  value       = aws_vpc.vpc.id
  description = "VPC ID"
}
output "public_subnets" {
    description = "VPC Public Subnets"
    value = [for az, subnet in aws_subnet.public_subnet: subnet.id]
}
output "private_subnets" {
    description = "VPC Private Subnets"
    value = [for az, subnet in aws_subnet.private_subnet: subnet.id]
}
output "eks_private_subnets" {
    description = "VPC Private Subnets - EKS"
    value = [for az, subnet in aws_subnet.eks_private_subnet: subnet.id]
}
output "vpc_cidr" {
    description = "VPC CIDR Block"
    value = aws_vpc.vpc.cidr_block
}
# Individual Subnets
# output "public_subnet_1" {
#     value = aws_subnet.public_subnet[0].id
#     description = "Public-Subnet-1 ID"
# }
# output "public_subnet_2" {
#     value = aws_subnet.public_subnet[1].id
#     description = "Public-Subnet-2 ID"
# }
# output "private_subnet_1" {
#     value = aws_subnet.private_subnet[0].id
#     description = "Private-Subnet-1 ID"
# }
# output "private_subnet_2" {
#     value = aws_subnet.private_subnet[1].id
#     description = "Private-Subnet-2 ID"
# }
# output "eks_private_subnet_1" {
#     value = aws_subnet.eks_private_subnet[0].id
#     description = "EKS-Private-Subnet-1 ID"
# }
# output "eks_private_subnet_2" {
#     value = aws_subnet.eks_private_subnet[1].id
#     description = "EKS-Private-Subnet-2 ID"
# }