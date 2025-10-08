terraform {
  required_version = "<version>"
  required_providers {
    pro {
      version = "<version-constraint>"
      source  = "<provider-address>"
    }
  }
  provider_meta "<LABEL>" {
    # Shown for completeness but only used for specific cases
  }
  backend "<TYPE>" {
    # `backend` is mutually exclusive with `cloud`
    path = "test"
  }
  cloud {
    # `cloud` is mutually exclusive with `backend`
    organization = "<organization-name>"
    workspaces {
      tags    = ["<tag>"]
      project = "<project-name>"
    }
    hostname = "app.terraform.io"

  }
  experiments = ["<feature-name>"]
}
