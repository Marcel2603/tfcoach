# core.enforce_parameter_order

Enforce the ordering of parameters as recommended by
the [Terraform docs](https://developer.hashicorp.com/terraform/language/style#resource-order):

1. If present, the `count` or `for_each` meta-argument
2. Resource-specific _non-block_ parameters
3. Resource-specific _block_ parameters
4. If required, a `lifecycle` block
5. If required, the `depends_on` parameter

This rule applies to the following blocks:

- `resource`
- `data`
- `module`
- `ephemeral`
- `output`

## Why

Using consistent and predictable ordering of parameters reduces the cognitive complexity and can lead to smaller diffs.

## Triggers

- Any `resource` block for which the parameters order differs from the recommendation

## Example

### Bad

```hcl
resource "aws_instance" "web1" {
  count = 1
  ami = 4321
  # depends_on should always come last
  depends_on = [
    aws_s3_bucket.s3
  ]
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_instance" "web2" {
  ami = 1234
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
  # this non-block parameter should come before the block parameter "instance_market_options"
  availability_zone = "custom-az"
}
```

### Good

```hcl
resource "aws_instance" "web1" {
  count = 1
  ami   = 4321
  lifecycle {
    ignore_changes = [tags]
  }
  depends_on = [
    aws_s3_bucket.s3
  ]
}

resource "aws_instance" "web2" {
  ami               = 1234
  availability_zone = "custom-az"
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
}
```

## Configuration

There is currently no configuration flags for that rule, besides the option to enable or disable the rule
