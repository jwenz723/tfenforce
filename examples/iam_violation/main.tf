resource "aws_iam_role" "service-role-default" {
  name = "test-default"
  description = "The default role that will used by kube2iam in eks clusters."
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "service-policy-default" {
  name        = "test-default"
  description = "This policy contains default permissions that will be granted to pods in eks clusters if a kube2iam role annotation is not specified on a pod."
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "*",
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "service-policy-attachment-default" {
  role = aws_iam_role.service-role-default.name
  policy_arn = aws_iam_policy.service-policy-default.arn
}