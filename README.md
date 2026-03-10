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
  05-locals.tf
  10-data.tf
  ...
  99-outputs.tf
```

Some files are expected to contain some block types exclusively:
* `00-variables.tf`: only `variable` blocks
* `01-terraform.tf`: only the `terraform` block with version & provider constraints
* `05-locals.tf`: only `locals` blocks
* `10-data.tf`: only `data` blocks
* `99-outputs.tf`: only `outputs` blocks

## Formatting

* Blocks must be spaced by a single empty line
* There must be a single empty line at the end of every file

We also recommend enforcing formatting via `terraform fmt`.

## Rules & configuration

All rules are enabled by default.

### `terraform_style_block_spacing`

```hcl
rule "terraform_style_block_spacing" {
  enabled  = true
}
```

Enforces blocks to be spaced by exactly one empty line.

### `terraform_style_file_end_newline`

```hcl
rule "terraform_style_file_end_newline" {
  enabled  = true
}
```

Enforces that every file ends with exactly one newline.

### `terraform_style_filename_pattern`

```hcl
rule "terraform_style_filename_pattern" {
  enabled  = true
}
```

Enforces that all Terraform files follow the `XX-name.tf` pattern (e.g., `00-variables.tf`).

### `terraform_style_variables_file`

```hcl
rule "terraform_style_variables_file" {
  enabled  = true
  filename = "00-variables"
}
```

Enforces all `variable` blocks to be in `00-variables.tf`, and that the file contains only `variable` blocks. The file name can be configured via the `filename` argument.

### `terraform_style_terraform_file`

```hcl
rule "terraform_style_terraform_file" {
  enabled  = true
  filename = "01-terraform"
}
```

Enforces the `terraform` block (version and provider constraints) to be in `01-terraform.tf`, that the file contains only the `terraform` block, and that there is exactly one such block. The file name can be configured via the `filename` argument.

### `terraform_style_locals_file`

```hcl
rule "terraform_style_locals_file" {
  enabled  = true
  filename = "05-locals"
}
```

Enforces all `locals` blocks to be in `05-locals.tf`, that the file contains only `locals` blocks, and that there is exactly one such block. The file name can be configured via the `filename` argument.

### `terraform_style_data_file`

```hcl
rule "terraform_style_data_file" {
  enabled  = true
  filename = "10-data"
}
```

Enforces all `data` blocks to be in `10-data.tf`, and that the file contains only `data` blocks. The file name can be configured via the `filename` argument.

### `terraform_style_outputs_file`

```hcl
rule "terraform_style_outputs_file" {
  enabled  = true
  filename = "99-outputs"
}
```

Enforces all `output` blocks to be in `99-outputs.tf`, and that the file contains only `output` blocks. The file name can be configured via the `filename` argument.

## Recommended configuration

We recommend the following `.tflint.hcl` configuration:

```hcl
plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

rule "terraform_documented_variables" { enabled = true }
rule "terraform_documented_outputs" { enabled = true }
rule "terraform_naming_convention" { enabled = true }
rule "terraform_comment_syntax" { enabled = true }
rule "terraform_unused_required_providers" { enabled = true }

plugin "terraform_style" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/Heldroe/tflint-ruleset-terraform-style"
}
```

## License

This project is licensed under the [MIT License](LICENSE).
