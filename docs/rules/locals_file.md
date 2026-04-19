# terraform_style_locals_file

Enforces that all `locals` blocks are consolidated into a designated file and that this file contains only authorized block types.

## Why?

Grouping local value definitions in a single file (usually `02-locals.tf`) prevents them from being scattered across multiple files, which can make it difficult to track down where a local variable is defined.

## Details

By default, this rule ensures that:
1. All `locals` blocks are defined in `02-locals.tf`.
2. `02-locals.tf` contains only `locals` blocks.
3. There is at most one `locals` block in the file.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `filename` | The base name of the file (without `.tf` extension). | string | `02-locals` |
| `allowed_blocks` | A list of block types allowed in the designated file. | list(string) | `["locals"]` |
| `exempt_blocks` | A map of block types to lists of subtypes that are exempt from the rule, even if their block type is not in `allowed_blocks`. | map(list(string)) | `{}` |

## Examples

### Denied

```hcl
# In 10-main.tf
locals { env = "prod" } # locals must be in 02-locals.tf

# In 02-locals.tf
resource "aws_s3_bucket" "b" { ... } # Resources are not allowed here
```

### Accepted

```hcl
# In 02-locals.tf
locals {
  name_prefix = "${var.project}-${var.env}"
  common_tags = {
    Project = var.project
  }
}
```
