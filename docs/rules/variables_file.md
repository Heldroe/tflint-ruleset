# terraform_style_variables_file

Enforces that `variable` blocks are placed in a designated file and that this file contains only authorized block types.

## Why?

Keeping variables in a dedicated file (usually `00-variables.tf`) makes it easier to understand the module's interface at a glance and maintains a clean separation of concerns.

## Details

By default, this rule ensures that:
1. All `variable` blocks are defined in `00-variables.tf`.
2. `00-variables.tf` contains only `variable` and `check` blocks.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `filename` | The base name of the file (without `.tf` extension). | string | `00-variables` |
| `allowed_blocks` | A list of block types allowed in the designated file. | list(string) | `["variable", "check"]` |
| `exempt_blocks` | A map of block types to lists of subtypes that are exempt from the rule, even if their block type is not in `allowed_blocks`. | map(list(string)) | `{}` |

## Examples

### Denied

```hcl
# In 10-main.tf
variable "instance_type" { ... } # Variables must be in 00-variables.tf

# In 00-variables.tf
resource "aws_instance" "foo" { ... } # Resources are not allowed here
```

### Accepted

```hcl
# In 00-variables.tf
variable "region" {
  type = string
}

check "health" {
  # ...
}
```
