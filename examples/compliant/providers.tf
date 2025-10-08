terraform {
  required_providers {
    aws = {
      version = "6.0"
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {}
