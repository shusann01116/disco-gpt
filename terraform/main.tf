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
