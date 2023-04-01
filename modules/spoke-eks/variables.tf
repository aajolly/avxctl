# locals = {
#   public_subnet_tags = {
#     "kubernetes.io/cluster/${var.cluster_name}" = "shared"
#     "kubernetes.io/role/elb"                    = 1
#   }
# 	private_subnet_tags = {
#     "kubernetes.io/cluster/${var.cluster_name}" = "shared"
#     "kubernetes.io/role/internal-elb"             = 1
#   }
# }
variable "enable_aws_load_balancer_controller" {
  type = bool
}
variable "cluster_name" {
  type = string
}
variable "eks_private_subnet_ids" {
  type = list
}
variable "eks_public_subnet_ids" {
  type = list
}