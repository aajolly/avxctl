data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "avx_role_ec2_policy" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    resources = ["arn:aws:iam::*:role/avx*"]
  }

  statement {
    actions = [
      "aws-marketplace:MeterUsage",
      "s3:GetBucketLocation",
    ]

    resources = ["*"]
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
  url = "https://ipapi.co/ip/"

  # Optional request headers
  request_headers = {
    Accept = "text/*"
  }
}

locals {
  account_id      = data.aws_caller_identity.current.account_id
  tool_prefix     = "avxctl"
  my_public_ip    = "${data.http.my_ip.response_body}/32"
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
}
variable "ctrl_version" {
  type    = string
  default = "latest"
}
