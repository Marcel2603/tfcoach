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
