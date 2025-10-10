# core.naming_convention

Enforce the naming convention.

## Why

Consistent naming improves module reuse and keeps downstream references simple.

## Triggers

- Any block whose not following the naming convention `a-z0-9_`.

## Example

### Bad

```hcl
resource "aws_s3_bucket" "Foo" {}
```

### Good

```hcl
resource "aws_s3_bucket" "foo" {}
```

## Configuration

There is currently no configuration flags for that rule, beside the option to enable or disable the rule
