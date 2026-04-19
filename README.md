# TFLint Ruleset for Terraform coding style

This TFLint ruleset enforces a consistent file structure and coding style for Terraform projects. It encourages organizing resources, variables, outputs, and locals into specific files and ensures proper formatting.

## Installation

To use this plugin, you can declare it in your `.tflint.hcl` file:

```hcl
plugin "terraform_style" {
  enabled = true
  version = "0.2.2"
  source  = "github.com/Heldroe/tflint-ruleset-terraform-style"
}
```

Then run `tflint --init`.

## Expected file names and content

All file names are expected to follow the `XX-name.tf` format to ensure a consistent ordering in file browsers.

```
module/
  00-variables.tf
  01-terraform.tf
  02-locals.tf
  03-data.tf
  ...
  10-ec2.tf
  20-helm.tf
  30-hetzner.tf
  ...
  99-outputs.tf
```

Some files are expected to contain specific block types exclusively:
* `00-variables.tf`: only `variable` and `check` blocks
* `01-terraform.tf`: only the `terraform` block with version & provider constraints
* `02-locals.tf`: only `locals` blocks
* `03-data.tf`: only `data` blocks
* `99-outputs.tf`: only `output` blocks

Resource files (all other `.tf` files) may only contain `check`, `module`, `moved`, `removed`, and `resource` blocks by default.

We recommend naming your files with `resource` blocks starting with `10-` onward, for example `10-s3.tf`.

All file rules support an `allowed_blocks` parameter to customize which block types are permitted, and an `exempt_blocks` parameter to allow specific subtypes of otherwise-disallowed blocks (e.g., allow `data "aws_iam_policy_document"` in resource files without allowing all `data` blocks).

## Formatting

* Blocks must be spaced by a single empty line
* There must be a single empty line at the end of every file

We also recommend enforcing formatting via `terraform fmt` and validation via `terraform validate`.

## Rules & configuration

All rules are enabled by default.

### [`terraform_style_block_spacing`](./docs/rules/block_spacing.md)

```hcl
rule "terraform_style_block_spacing" {
  enabled  = true
}
```

Enforces blocks to be spaced by exactly one empty line.

### [`terraform_style_file_end_newline`](./docs/rules/file_end_newline.md)

```hcl
rule "terraform_style_file_end_newline" {
  enabled  = true
}
```

Enforces that every file ends with exactly one newline.

### [`terraform_style_filename_pattern`](./docs/rules/filename_pattern.md)

```hcl
rule "terraform_style_filename_pattern" {
  enabled   = true
  min_index = 10
}
```

Enforces that all Terraform files follow the `XX-name.tf` pattern (e.g., `00-variables.tf`). Resource files must use an index equal or greater than `min_index` (default `10`).

### [`terraform_style_variables_file`](./docs/rules/variables_file.md)

```hcl
rule "terraform_style_variables_file" {
  enabled        = true
  filename       = "00-variables"
  allowed_blocks = ["variable", "check"]
  exempt_blocks  = {}
}
```

Enforces that `00-variables.tf` contains only the allowed block types. By default, `variable` and `check` blocks are allowed. The file name, allowed blocks, and block exemptions can be configured.

### [`terraform_style_terraform_file`](./docs/rules/terraform_file.md)

```hcl
rule "terraform_style_terraform_file" {
  enabled        = true
  filename       = "01-terraform"
  allowed_blocks = ["terraform"]
  exempt_blocks  = {}
}
```

Enforces that `01-terraform.tf` contains only the allowed block types and at most one `terraform` block. By default, only the `terraform` block is allowed. The file name, allowed blocks, and block exemptions can be configured.

### [`terraform_style_locals_file`](./docs/rules/locals_file.md)

```hcl
rule "terraform_style_locals_file" {
  enabled        = true
  filename       = "02-locals"
  allowed_blocks = ["locals"]
  exempt_blocks  = {}
}
```

Enforces that `02-locals.tf` contains only the allowed block types and at most one `locals` block. By default, only `locals` blocks are allowed. The file name, allowed blocks, and block exemptions can be configured.

### [`terraform_style_data_file`](./docs/rules/data_file.md)

```hcl
rule "terraform_style_data_file" {
  enabled        = true
  filename       = "03-data"
  allowed_blocks = ["data"]
  exempt_blocks  = {}
}
```

Enforces that `03-data.tf` contains only the allowed block types. By default, only `data` blocks are allowed. The file name, allowed blocks, and block exemptions can be configured.

### [`terraform_style_outputs_file`](./docs/rules/outputs_file.md)

```hcl
rule "terraform_style_outputs_file" {
  enabled        = true
  filename       = "99-outputs"
  allowed_blocks = ["output"]
  exempt_blocks  = {}
}
```

Enforces that `99-outputs.tf` contains only the allowed block types. By default, only `output` blocks are allowed. The file name, allowed blocks, and block exemptions can be configured.

### [`terraform_style_resource_file`](./docs/rules/resource_file.md)

```hcl
rule "terraform_style_resource_file" {
  enabled        = true
  allowed_blocks = ["check", "module", "moved", "removed", "resource"]
  exempt_blocks  = {}
}
```

Enforces that resource files contain only the allowed block types. By default, `data`, `import`, `locals`, `output`, `provider`, `terraform`, and `variable` blocks are not allowed in resource files. Use `exempt_blocks` to allow specific subtypes of otherwise-disallowed blocks.

### [`terraform_style_no_backend_block`](./docs/rules/no_backend_block.md)

```hcl
rule "terraform_style_no_backend_block" {
  enabled = true
}
```

Enforces that no `backend` or `cloud` blocks are defined in the `terraform` block. Backend and cloud configuration should be passed via CLI or defined in a separate file if using a partial configuration.

### [`terraform_style_no_empty_file`](./docs/rules/no_empty_file.md)

```hcl
rule "terraform_style_no_empty_file" {
  enabled = true
}
```

Enforces that Terraform files are not empty (must define at least one block or attribute).

### [`terraform_style_trailing_comma`](./docs/rules/trailing_comma.md)

```hcl
rule "terraform_style_trailing_comma" {
  enabled = true
  exclude_single_element = false
}
```

Enforces trailing comma rules for multi-line lists and maps.

### [`terraform_style_map_assignment`](./docs/rules/map_assignment.md)

```hcl
rule "terraform_style_map_assignment" {
  enabled = true
}
```

Enforces that maps/objects use the equal sign `=` for key-value assignment instead of the colon `:`.

### [`terraform_style_comment_style`](./docs/rules/comment_style.md)

```hcl
rule "terraform_style_comment_style" {
  enabled = true
}
```

Enforces consistent comment style (only `#` allowed, single space required).

### [`terraform_style_resource_arguments`](./docs/rules/resource_arguments.md)

```hcl
rule "terraform_style_resource_arguments" {
  enabled = true
}
```

Enforces ordering and spacing of standard arguments and blocks (count, for_each, depends_on, lifecycle, etc).

### [`terraform_style_variable_arguments`](./docs/rules/variable_arguments.md)

```hcl
rule "terraform_style_variable_arguments" {
  enabled = true
  order   = ["type", "nullable", "sensitive", "ephemeral", "default", "description", "validation"]
}
```

Enforces a consistent ordering of arguments within `variable` blocks. Arguments not listed in `order` are ignored. Only present arguments are checked — none are required to be present. The order can be customized via the `order` parameter.

### [`terraform_style_output_arguments`](./docs/rules/output_arguments.md)

```hcl
rule "terraform_style_output_arguments" {
  enabled = true
  order   = ["description", "sensitive", "ephemeral", "value", "precondition", "depends_on"]
}
```

Enforces a consistent ordering of arguments within `output` blocks. Arguments not listed in `order` are ignored. Only present arguments are checked — none are required to be present. The order can be customized via the `order` parameter.

### [`terraform_style_block_internal_spacing`](./docs/rules/block_internal_spacing.md)

```hcl
rule "terraform_style_block_internal_spacing" {
  enabled = true
}
```

Enforces consistent spacing inside resource, module, and other blocks.

### [`terraform_style_structure_layout`](./docs/rules/structure_layout.md)

```hcl
rule "terraform_style_structure_layout" {
  enabled = true
}
```

Enforces consistent layout for lists and maps (closing bracket placement, element on new line).

### [`terraform_style_no_provider_argument`](./docs/rules/no_provider_argument.md)

```hcl
rule "terraform_style_no_provider_argument" {
  enabled = false
}
```

Forbids the use of the `provider` argument within resource and data blocks. **Disabled by default.**

## Recommended configuration

We recommend the following `.tflint.hcl` configuration:

```hcl
# Terraform plugin for TFLint
plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

rule "terraform_documented_variables" {
  enabled = true
}

rule "terraform_documented_outputs" {
  enabled = true
}

rule "terraform_naming_convention" {
  enabled = true
}

rule "terraform_comment_syntax" {
  enabled = true
}

rule "terraform_unused_required_providers" {
  enabled = true
}

# Terraform style plugin
plugin "terraform_style" {
  enabled = true
  version = "0.2.2"
  source  = "github.com/Heldroe/tflint-ruleset-terraform-style"
}
```

When using pure Terraform, you might want to adjust the following for your top level modules:
- Add `"provider"` to `allowed_blocks` in `terraform_style_resource_file` (or whichever file rule should contain provider blocks)
- Disable `terraform_style_no_backend_block`

## Contributing

If you want to contribute to this project, please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).
