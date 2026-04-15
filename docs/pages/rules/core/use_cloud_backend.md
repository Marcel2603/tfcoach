# core.use_cloud_backend

This rule should enforce that no local statefile is being used.

## Why

If you store the state locally (default backend), it will be impossible to run terraform from other PCs

## Triggers

- No `backend` or `cloud` block inside `terraform`
- `backend`-block of type `local`

## Example

### Bad

```hcl
terraform {
 backend "local" {
 }
}
```

or no backend at all

### Good

```hcl
terraform {
  backend "s3" {
    bucket = "mybucket"
    key    = "path/to/my/key"
    region = "us-east-1"
  }
}
```

## Configuration

There is currently no configuration flags for that rule, beside the option to enable or disable the rule
