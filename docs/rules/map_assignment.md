# terraform_style_map_assignment

Enforces that maps/objects use the equal sign `=` for key-value assignment instead of the colon `:`.

## Why?

Using `=` consistently with attributes makes Terraform code more uniform and easier to read.

## Examples

### Denied

```hcl
tags = {
  Environment : "prod"
}
```

### Accepted

```hcl
tags = {
  Environment = "prod"
}
```
