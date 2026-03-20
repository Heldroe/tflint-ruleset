# terraform_style_resource_file

Enforces that resource files (any `.tf` file not covered by other specific file rules) only contain authorized block types.

## Why?

By restricting what can be placed in resource files, this rule ensures that core module configuration like variables, outputs, and locals doesn't get buried among resource definitions.

## Details

By default, any Terraform file that matches the `XX-name.tf` pattern and is NOT one of the standard files (variables, terraform, locals, data, outputs) is considered a resource file.

Default allowed block types: `check`, `module`, `moved`, `removed`, `resource`.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `allowed_blocks` | A list of block types allowed in resource files. | list(string) | `["check", "module", "moved", "removed", "resource"]` |

## Examples

### Denied

```hcl
# In 10-network.tf
variable "vpc_cidr" { ... } # variables are not allowed in resource files

locals { ... } # locals are not allowed in resource files
```

### Accepted

```hcl
# In 10-network.tf
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

module "subnets" {
  source = "./modules/subnets"
  # ...
}
```
