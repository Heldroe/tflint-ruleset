# terraform_style_no_empty_file

Enforces that Terraform files are not empty. Every file must define at least one block or attribute.

## Why?

Empty files are unnecessary and can clutter the project structure. If a file is not needed, it should be deleted.

## Examples

### Denied

```hcl
# main.tf is empty
```

### Accepted

```hcl
resource "null_resource" "foo" {}
```
