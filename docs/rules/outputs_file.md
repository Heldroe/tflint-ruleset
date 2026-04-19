# terraform_style_outputs_file

Enforces that all `output` blocks are consolidated into a designated file and that this file contains only authorized block types.

## Why?

Keeping outputs in a single file (usually `99-outputs.tf`) makes it easy to understand what information the module exposes to its users.

## Details

By default, this rule ensures that:
1. All `output` blocks are defined in `99-outputs.tf`.
2. `99-outputs.tf` contains only `output` blocks.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `filename` | The base name of the file (without `.tf` extension). | string | `99-outputs` |
| `allowed_blocks` | A list of block types allowed in the designated file. | list(string) | `["output"]` |
| `exempt_blocks` | A map of block types to lists of subtypes that are exempt from the rule, even if their block type is not in `allowed_blocks`. | map(list(string)) | `{}` |

## Examples

### Denied

```hcl
# In 10-main.tf
output "vpc_id" { ... } # outputs must be in 99-outputs.tf

# In 99-outputs.tf
resource "aws_db_instance" "db" { ... } # Resources are not allowed here
```

### Accepted

```hcl
# In 99-outputs.tf
output "instance_ip" {
  value = aws_instance.web.public_ip
}
```
