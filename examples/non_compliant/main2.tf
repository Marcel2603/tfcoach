data "archive_file" "zip" {}

locals {}

resource "null_resource" "tEst" {}

resource "test" "is-complaint" {}

output "test2" {
  value = "test"
}

provider "aws2" {}

terraform {}

variable "test2" {}


variable "environment" {
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be dev, staging, or prod."
  }

  type        = string
}
