# tfcoach-ignore: rule1, rule2
terraform {
  required_version = "1.0"
  experiments = ["<feature-name>"]
  required_providers {
    aws = {
      version = "6.0"
      source  = "hashicorp/aws"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.7.1"
    }
    null = {
      source  = "mildred/null"
      version = "1.1.0"
    }
  }
}
