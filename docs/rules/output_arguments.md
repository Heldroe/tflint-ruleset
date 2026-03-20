# terraform_style_output_arguments

Enforces a consistent ordering of arguments within `output` blocks.

## Why?

A consistent ordering of arguments makes it easier to scan output definitions and quickly find specific properties like the value or description.

## Details

This rule checks the order of arguments present in an `output` block. It does NOT require any arguments to be present; it only checks the relative order of those that are.

Arguments not specified in the `order` configuration are ignored.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `order` | The desired order of arguments. | list(string) | `["description", "sensitive", "ephemeral", "value", "precondition", "depends_on"]` |

## Examples

### Denied

```hcl
output "foo" {
  value       = "bar"
  description = "A foo" # 'description' should come before 'value'
}
```

### Accepted

```hcl
output "foo" {
  description = "A foo"
  value       = "bar"
}
```
