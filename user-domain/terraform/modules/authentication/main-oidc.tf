provider "aws" {
  region = "ap-southeast-1"
}

resource "aws_iam_openid_connect_provider" "github_action" {
  client_id_list = ["sts.amazonaws.com"]
  url = "https://token.actions.githubusercontent.com"
  thumbprint_list = [data.tls_certificate.github.certificates[0].sha1_fingerprint]
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect = "Allow"

    principals {
      identifiers = [aws_iam_openid_connect_provider.github_action.arn]
      type = "Federated"
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:sub"
      values = [
        for a in var.allowed_repos_branches : "repo:${a["org"]}/${a["repo"]}:ref:refs/heads/${a["branch"]}"
      ]
    }
  }
}

resource "aws_iam_role" "auth_role" {
  name = "AuthenticationRole"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "auth_role_admin" {
  role = aws_iam_role.auth_role.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}