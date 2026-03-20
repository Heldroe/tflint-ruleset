# terraform_style_variable_arguments

Enforces a consistent ordering of arguments within `variable` blocks.

## Why?

A consistent ordering of arguments makes it easier to scan variable definitions and quickly find specific properties like the type or default value.

## Details

This rule checks the order of arguments present in a `variable` block. It does NOT require any arguments to be present; it only checks the relative order of those that are.

Arguments not specified in the `order` configuration are ignored.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `order` | The desired order of arguments. | list(string) | `["type", "nullable", "sensitive", "ephemeral", "default", "description", "validation"]` |

## Examples

### Denied

```hcl
variable "foo" {
  description = "A foo"
  type        = string # 'type' should come before 'description'
}
```

### Accepted

```hcl
variable "foo" {
  type        = string
  description = "A foo"
}
```
