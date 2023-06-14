# resource "aws_iam_role" "aajolly_test_role" {
#   name        = "aajolly_test_role"
#   path        = "/"
#   description = "Test IAM Role"
#   managed_policy_arns = [
#     "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
#   ]
#   assume_role_policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Action = "sts:AssumeRole"
#         Effect = "Allow"
#         Sid    = ""
#         Principal = {
#           Service = "ec2.amazonaws.com"
#         }
#       },
#     ]
#   })

#   tags = {
#     Name = "aajolly_test_role"
#   }
# }

# resource "aws_iam_instance_profile" "aajolly_test_profile" {
#   name = "aajolly_test_profile"
#   role = aws_iam_role.aajolly_test_role.name
# }