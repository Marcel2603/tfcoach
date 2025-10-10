# core.required_provider_must_be_declared

Enforce that all providers used in resources or data sources are declared in the `terraform.required_providers` block.

## Why

Using non explicitly declared providers usually leads to bugs when trying to apply changes.

## Triggers

- Any `resource` or `data` block for which the provider is not declared

## Example

### Bad

```hcl
resource "aws_s3_bucket" "this" {}
```

### Good

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_s3_bucket" "this" {}
```

Note that the provider declaration and the resource usage do not need to be in the same file.

## Configuration

There is currently no configuration flags for that rule, beside the option to enable or disable the rule
