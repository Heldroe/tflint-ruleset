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
| `exempt_blocks` | A map of block types to lists of subtypes that are exempt from the rule, even if their block type is not in `allowed_blocks`. | map(list(string)) | `{}` |

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

### Using `exempt_blocks`

```hcl
# .tflint.hcl
rule "terraform_style_resource_file" {
  enabled = true
  exempt_blocks = {
    "data" = ["aws_iam_policy_document"]
  }
}
```

With the configuration above, `data "aws_iam_policy_document"` blocks are allowed in resource files, while other `data` blocks (e.g., `data "aws_ami"`) are still flagged.
