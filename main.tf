terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_api_gateway_rest_api" "release_api" {
  name          = "release-api"
}

resource "aws_dynamodb_table" "release_table" {
  name           = "release-table"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "ClientId"

  attribute {
    name = "ClientId"
    type = "S"
  }

}

data "aws_iam_policy_document" "get_flags_policy" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:Query*"]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "get_flags_policy" {
  name        = "get-flags-policy"
  description = "policy for get flags role"
  policy      = data.aws_iam_policy_document.get_flags_policy.json
}

data "aws_iam_policy_document" "assume_role" {
    statement {
      effect = "Allow"

      principals {
         type = "Service"
         identifiers = ["lambda.amazonaws.com"]
      }

      actions = ["sts:AssumeRole"]
    }
}

resource "aws_iam_role" "get_flags_role" {
  name = "get-flags-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "get_flags_attach_role_policy" {
  role       = aws_iam_role.get_flags_role.name
  policy_arn = aws_iam_policy.get_flags_policy.arn
}

resource "aws_lambda_function" "get_flags_lambda" {
  filename      = "src/getFlags/getFlags.zip"
  function_name = "getFlags"
  handler       = "getFlags"
  role = aws_iam_role.get_flags_role.arn
  source_code_hash = filebase64sha256("src/getFlags/getFlags.zip")
  runtime = "go1.x"
  environment {
    variables = {
      RELEASE_TABLE_ARN = aws_dynamodb_table.release_table.arn
    }
  }
}

resource "aws_api_gateway_resource" "release_api_resource" {
  rest_api_id = aws_api_gateway_rest_api.release_api.id
  parent_id = aws_api_gateway_rest_api.release_api.root_resource_id
  path_part = "release"
}

resource "aws_api_gateway_method" "release_api_get_method" {
  rest_api_id = aws_api_gateway_rest_api.release_api.id
  resource_id = aws_api_gateway_resource.release_api_resource.id
  http_method = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "release_get_integration" {
    rest_api_id = aws_api_gateway_rest_api.release_api.id
    resource_id = aws_api_gateway_resource.release_api_resource.id
    http_method = aws_api_gateway_method.release_api_get_method.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.get_flags_lambda.invoke_arn
}

resource "aws_lambda_permission" "release_lambda_permission" {
  statement_id  = "AllowReleaseInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_flags_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_api_gateway_rest_api.release_api.execution_arn}/*"
}