# terraform_style_data_file

Enforces that all `data` blocks are consolidated into a designated file and that this file contains only authorized block types.

## Why?

Consolidating data sources into a single file (usually `03-data.tf`) makes it easy to see all external dependencies and information being pulled into the module.

## Details

By default, this rule ensures that:
1. All `data` blocks are defined in `03-data.tf`.
2. `03-data.tf` contains only `data` blocks.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `filename` | The base name of the file (without `.tf` extension). | string | `03-data` |
| `allowed_blocks` | A list of block types allowed in the designated file. | list(string) | `["data"]` |
| `exempt_blocks` | A map of block types to lists of subtypes that are exempt from the rule, even if their block type is not in `allowed_blocks`. | map(list(string)) | `{}` |

## Examples

### Denied

```hcl
# In 10-main.tf
data "aws_ami" "ubuntu" { ... } # data blocks must be in 03-data.tf

# In 03-data.tf
resource "aws_security_group" "sg" { ... } # Resources are not allowed here
```

### Accepted

```hcl
# In 03-data.tf
data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
```
