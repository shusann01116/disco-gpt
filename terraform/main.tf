provider "aws" {
  region = "ap-northeast-1"

  # Make it faster by skipping something
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_credentials_validation = true

  default_tags {
    tags = {
      Repo = "github.com/shusann01116/disco-gpt"
    }
  }
}

resource "random_id" "this" {
  byte_length = 8
}

# HTTP API Gateway
module "api_gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "~> 2.0"

  name          = "disco-gpt-${random_id.this.id}"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  domain_name                 = "${var.subdomain_name}.${var.domain_name}"
  domain_name_certificate_arn = module.acm.acm_certificate_arn

  integrations = {
    "GET /" = {
      lambda_arn             = module.lambda.lambda_function_arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 12000
    }

    "POST /" = {
      lambda_arn             = module.lambda.lambda_function_arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 12000
    }
  }
}

# Lambda

module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~>4.17.0"

  function_name = "disco-gpt-${random_id.this.id}"
  handler       = "main"
  runtime       = "go1.x"

  publish = true

  create_package         = false
  local_existing_package = data.archive_file.lambda_package.output_path

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*/*"
    }
  }

  environment_variables = {
    "DISCORD_PUBLIC_KEY" = var.discord_public_key
  }
}

resource "null_resource" "build" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 go build -o ../bin/main ../main.go"
  }
}

data "archive_file" "lambda_package" {
  depends_on  = [null_resource.build]
  type        = "zip"
  source_file = "../bin/main"
  output_path = "../bin/main.zip"
}

# ACM

data "aws_route53_zone" "this" {
  name = var.domain_name
}

module "acm" {
  source  = "terraform-aws-modules/acm/aws"
  version = "~> 3.0"

  domain_name               = var.domain_name
  zone_id                   = data.aws_route53_zone.this.zone_id
  subject_alternative_names = ["${var.subdomain_name}.${var.domain_name}"]
}

# Route53

resource "aws_route53_record" "api" {
  zone_id = data.aws_route53_zone.this.zone_id
  name    = var.subdomain_name
  type    = "A"

  alias {
    name                   = module.api_gateway.apigatewayv2_domain_name_configuration[0].target_domain_name
    zone_id                = module.api_gateway.apigatewayv2_domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}
