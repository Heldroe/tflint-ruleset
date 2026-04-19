# terraform_style_terraform_file

Enforces that the `terraform` configuration block is placed in a designated file and that this file contains only authorized block types.

## Why?

Consolidating Terraform version and provider constraints into a single, predictable file (usually `01-terraform.tf`) ensures that module requirements are easy to locate and maintain.

## Details

By default, this rule ensures that:
1. The `terraform` block is defined in `01-terraform.tf`.
2. `01-terraform.tf` contains only the `terraform` block.
3. There is at most one `terraform` block in the file.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `filename` | The base name of the file (without `.tf` extension). | string | `01-terraform` |
| `allowed_blocks` | A list of block types allowed in the designated file. | list(string) | `["terraform"]` |
| `exempt_blocks` | A map of block types to lists of subtypes that are exempt from the rule, even if their block type is not in `allowed_blocks`. | map(list(string)) | `{}` |

## Examples

### Denied

```hcl
# In 10-main.tf
terraform { ... } # terraform block must be in 01-terraform.tf

# In 01-terraform.tf
resource "aws_vpc" "main" { ... } # Resources are not allowed here
```

### Accepted

```hcl
# In 01-terraform.tf
terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```
