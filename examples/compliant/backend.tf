# tfcoach-ignore-file: rule1
terraform {
  backend "remote" {}
  cloud {
    workspaces {}
  }
  backend "remote" {}
  cloud {
    workspaces {}
  }
}
