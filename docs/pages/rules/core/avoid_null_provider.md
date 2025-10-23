# core.avoid_null_provider

Enforces that the [hashicorp/null](https://registry.terraform.io/providers/hashicorp/null/latest/docs)
 provider is not used.

## Why

The **hashicorp/null** provider was widely used in older Terraform versions for “glue logic,” but it’s now
considered obsolete and discouraged.

You should not use it anymore and replace it with `terraform_data` for `null_resource` and `locals` for `null_data_source`

## Triggers

- Any usage of `null_resource` or `null_data_source`

## Example

### Bad

```hcl
resource "null_resource" "config" {
  triggers = {
    config = var.config
  }
}
```

```hcl
data "null_data_source" "example" {
  inputs = {
    full_id = "${var.name}-${var.region}"
  }
}
```

### Good

```hcl
resource "terraform_data" "config" {
  input = var.config
}

locals {
  full_id = "${var.name}-${var.region}"
}
```

## Configuration

There is currently no configuration flags for that rule, beside the option to enable or disable the rule
