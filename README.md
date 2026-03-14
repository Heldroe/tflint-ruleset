# TFLint Ruleset for Terraform coding style

This TFLint ruleset enforces a consistent file structure and coding style for Terraform projects. It encourages organizing resources, variables, outputs, and locals into specific files and ensures proper formatting.

## Installation

To use this plugin, you can declare it in your `.tflint.hcl` file:

```hcl
plugin "terraform-style" {
  enabled = true
  version = "0.1.0"
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

Some files are expected to contain some block types exclusively:
* `00-variables.tf`: only `variable` blocks
* `01-terraform.tf`: only the `terraform` block with version & provider constraints
* `02-locals.tf`: only `locals` blocks
* `03-data.tf`: only `data` blocks
* `99-outputs.tf`: only `outputs` blocks

We recommend naming your files with `resource` blocks starting with `10-` onward, for example `10-s3.tf`.

## Formatting

* Blocks must be spaced by a single empty line
* There must be a single empty line at the end of every file

We also recommend enforcing formatting via `terraform fmt`.

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
  enabled  = true
  filename = "00-variables"
}
```

Enforces all `variable` blocks to be in `00-variables.tf`, and that the file contains only `variable` blocks. The file name can be configured via the `filename` argument.

### [`terraform_style_terraform_file`](./docs/rules/terraform_file.md)

```hcl
rule "terraform_style_terraform_file" {
  enabled  = true
  filename = "01-terraform"
}
```

Enforces the `terraform` block (version and provider constraints) to be in `01-terraform.tf`, that the file contains only the `terraform` block, and that there is exactly one such block. The file name can be configured via the `filename` argument.

### [`terraform_style_locals_file`](./docs/rules/locals_file.md)

```hcl
rule "terraform_style_locals_file" {
  enabled  = true
  filename = "02-locals"
}
```

Enforces all `locals` blocks to be in `02-locals.tf`, that the file contains only `locals` blocks, and that there is exactly one such block. The file name can be configured via the `filename` argument.

### [`terraform_style_data_file`](./docs/rules/data_file.md)

```hcl
rule "terraform_style_data_file" {
  enabled  = true
  filename = "03-data"
}
```

Enforces all `data` blocks to be in `03-data.tf`, and that the file contains only `data` blocks. The file name can be configured via the `filename` argument.

### [`terraform_style_outputs_file`](./docs/rules/outputs_file.md)

```hcl
rule "terraform_style_outputs_file" {
  enabled  = true
  filename = "99-outputs"
}
```

Enforces all `output` blocks to be in `99-outputs.tf`, and that the file contains only `output` blocks. The file name can be configured via the `filename` argument.

### [`terraform_style_no_provider_block`](./docs/rules/no_provider_block.md)

```hcl
rule "terraform_style_no_provider_block" {
  enabled = true
}
```

Enforces that no `provider` blocks are defined in the module. Provider configurations should be passed from the root module or defined outside the module.

### [`terraform_style_no_backend_block`](./docs/rules/no_backend_block.md)

```hcl
rule "terraform_style_no_backend_block" {
  enabled = true
}
```

Enforces that no `backend` blocks are defined in the `terraform` block. Backend configuration should be passed via CLI or defined in a separate file if using a partial configuration.

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
  version = "0.1.0"
  source  = "github.com/Heldroe/tflint-ruleset-terraform-style"
}
```

When using pure Terraform, you might want to disable the following rules for your top level modules:
- `terraform_style_no_provider_block`
- `terraform_style_no_backend_block`

## License

This project is licensed under the [MIT License](LICENSE).
