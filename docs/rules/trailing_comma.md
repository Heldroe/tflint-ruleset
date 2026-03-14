# terraform_style_trailing_comma

Enforces trailing comma rules for multi-line lists and maps.

## Details

- **Lists** defined in multiple lines **must** have a trailing comma on the last line.
- **Maps** defined in multiple lines **must not** have any trailing comma on any lines.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `exclude_single_element` | If true, do not require a trailing comma for single-element multi-line lists. | boolean | `false` |

## Examples

### Denied

```hcl
# List missing trailing comma
availability_zones = [
  "us-east-1a",
  "us-east-1b"
]

# Map with trailing comma
tags = {
  Environment = "prod",
}
```

### Accepted

```hcl
# List with trailing comma
availability_zones = [
  "us-east-1a",
  "us-east-1b",
]

# Map without trailing comma
tags = {
  Environment = "prod"
}

# Single-element list (valid if exclude_single_element = true)
security_groups = [
  {
    name = "allow_all"
  }
]
```
