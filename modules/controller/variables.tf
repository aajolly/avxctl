data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

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
  name_prefix     = "aviatrix"
  my_public_ip    = jsondecode(data.http.my_ip.response_body).ip
  images_byol     = jsondecode(data.http.avx_iam_id.response_body).BYOL
  images_platinum = jsondecode(data.http.avx_iam_id.response_body).MeteredPlatinum
  ami_id          = var.type == "BYOL" || var.type == "byol" ? local.images_byol[data.aws_region.current.name] : local.images_platinum[data.aws_region.current.name]
}

variable "type" {
  default     = "byol"
  type        = string
  description = "Type of billing, can be 'meteredplatinum' or 'BYOL'."
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
variable "ebs-encryption" {
  type    = bool
  default = true
}
variable "admin_password" {
  type      = string
  default   = "Pa$$w0rd123"
  sensitive = true
}
variable "admin_email" {
  type      = string
  default   = "ajolly@aviatrix.com"
  sensitive = true
}
variable "ctrl_name" {
  type    = string
  default = "aajolly-test-ctrl"
}
variable "ctrl_customer_id" {
  type      = string
  default   = "aviatrixlab.com-1639517259.96"
  sensitive = true
}
variable "ctrl_version" {
  type    = string
  default = "latest"
}