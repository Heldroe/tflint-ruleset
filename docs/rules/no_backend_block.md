# terraform_style_no_backend_block

Enforces that no `backend` or `cloud` blocks are defined within the `terraform` configuration block.

## Why?

Hardcoding backend or cloud configuration within a module makes it less reusable across different environments. It is recommended to pass these configurations via CLI or using partial configuration files during `terraform init`.

## Details

This rule scans all `terraform` blocks and emits an issue if a `backend` or `cloud` block is found inside.

## Examples

### Denied

```hcl
terraform {
  backend "s3" {
    bucket = "my-terraform-state"
    # ...
  }
}
```

### Accepted

```hcl
terraform {
  required_version = ">= 1.0"
}
```
