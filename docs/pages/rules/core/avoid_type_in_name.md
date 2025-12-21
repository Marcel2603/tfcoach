# core.avoid_type_in_name

Enforces that the block-type got not repeated in the block-name.

## Why

Resource names shouldn’t repeat their type; this causes redundancy.

## Triggers

- Any repatation of the type-definition inside the name

## Example

### Bad

```hcl

resource "aws_s3_bucket" "s3_bucket_example" {} # repeats s3 and bucket
resource "aws_s3_bucket" "aws_example" {} # repeats aws
resource "aws_s3_bucket" "test_bucket" {} # repeats bucket
resource "aws_s3_bucket" "s3_test" {} # repeats s3
```

### Good

```hcl
resource "aws_s3_bucket" "ui" {}
resource "aws_s3_bucket" "this" {}
```

## Configuration

There is currently no configuration flags for that rule, beside the option to enable or disable the rule
