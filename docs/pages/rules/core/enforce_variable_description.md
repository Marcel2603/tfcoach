# core.enforce_variable_description

Enforces that all declared variables have a non-empty description

## Why

Even though the intent of variables may seem trivial at first, the variable name itself usually does not carry
enough information to ensure it stays that way. A good description can save much debugging time.

## Triggers

- Any variables with no description or an empty description

## Example

### Bad

```hcl
variable "test" {
  type = string
}
```

```hcl
variable "test2" {
  type = string
  description = ""
}
```

### Good

```hcl
variable "test" {
  description = "some descriptive text"
  type = string
}
```

## Configuration

There is currently no configuration flags for that rule, beside the option to enable or disable the rule
