data "archive_file" "zip" {}

locals {}

resource "null_resource" "tEst" {}

resource "test_resource" "is-not-compliant" {}

output "test" {
  value = "test"
}

provider "aws" {}

terraform {}

variable "test" {}

resource "azurerm_resource_group" "this" {
  location = "test"
  name = "hello"
}

resource "aws_s3_bucket" "s3" {
  bucket = "non-compliant-naming"
}

resource "aws_instance" "web1" {
  count = 1
  ami = 4321
  depends_on = [
    aws_s3_bucket.s3
  ]
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_instance" "web2" {
  ami = 1234
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
  availability_zone = "custom-az"
}
