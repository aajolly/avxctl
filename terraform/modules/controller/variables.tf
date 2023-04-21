data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

data "http" "avx_role_ec2_policy" {
  url = "https://s3-us-west-2.amazonaws.com/aviatrix-download/iam_assume_role_policy.txt"
  request_headers = {
    "Accept" = "text/*"
  }
}
data "http" "avx_role_app_policy" {
  url = "https://s3-us-west-2.amazonaws.com/aviatrix-download/IAM_access_policy_for_CloudN.txt"
  request_headers = {
    "Accept" = "text/*"
  }
}

data "http" "avx_iam_id" {
  url = "https://s3-us-west-2.amazonaws.com/aviatrix-download/AMI_ID/ami_id.json"
  request_headers = {
    "Accept" = "application/json"
  }
}
data "http" "my_ip" {
  url    = "https://api.ipify.org?format=json"
  method = "GET"
}

locals {
  account_id      = data.aws_caller_identity.current.account_id
  tool_prefix     = "avxctl"
  my_public_ip    = "${jsondecode(data.http.my_ip.response_body).ip}/32"
  images_byol     = jsondecode(data.http.avx_iam_id.response_body).BYOL
  images_platinum = jsondecode(data.http.avx_iam_id.response_body).MeteredPlatinum
  ami_id          = var.type == "BYOL" || var.type == "byol" ? local.images_byol[data.aws_region.current.name] : local.images_platinum[data.aws_region.current.name]
  region          = coalesce(var.region, data.aws_region.current.name)
}

variable "type" {
  default     = "byol"
  type        = string
  description = "Type of billing, can be 'meteredplatinum' or 'BYOL'."
}
variable "region" {
  type = string
  default = ""
}
variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/24"
}
variable "internet_cidr" {
  type    = string
  default = "0.0.0.0/0"
}
variable "instance_type" {
  type    = string
  default = "t3.large"
}
variable "root_volume_size" {
  type    = number
  default = 64
}
variable "root_volume_type" {
  type    = string
  default = "gp2"
}
variable "ebs_encryption" {
  type    = bool
  default = true
}
variable "keypair" {
  type = string
  default = "aajolly-apse2"
}
# variable "admin_password" {
#   type      = string
#   default   = "Pa$$w0rd123"
#   sensitive = true
# }
# variable "admin_email" {
#   type      = string
#   default   = "ajolly@aviatrix.com"
#   sensitive = true
# }
# variable "ctrl_name" {
#   type    = string
#   default = "avxctl-controller"
# }
# variable "ctrl_customer_id" {
#   type      = string
#   sensitive = true
# }
# variable "ctrl_version" {
#   type    = string
#   default = "latest"
# }