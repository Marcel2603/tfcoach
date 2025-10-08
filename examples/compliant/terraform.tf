terraform {
<<<<<<< HEAD
  required_version = "1.0"
  experiments = ["<feature-name>"]
  required_providers {
    aws = {
      version = "6.0"
      source = "hashicorp/aws"
=======
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.7.1"
    }
    null = {
      source  = "mildred/null"
      version = "1.1.0"
>>>>>>> 5b7adb0 (fix: Actually call the new rule Finish method, extend examples)
    }
  }
}
