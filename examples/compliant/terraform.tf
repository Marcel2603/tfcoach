terraform {
  required_version = "1.0"
  experiments = ["<feature-name>"]
  required_providers {
    aws = {
      version = "6.0"
      source = "hashicorp/aws"
    }
  }
}
