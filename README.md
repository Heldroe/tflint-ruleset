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
  01-setup.tf
  05-locals.tf
  10-data.tf
  ...
  99-outputs.tf
```

Some files are expected to contain some block types exclusively:
* `00-variables.tf`: only `variable` blocks
* `01-setup.tf`: only the `terraform` block with version & provider constraints
* `05-locals.tf`: only `locals` blocks
* `10-data.tf`: only `data` blocks
* `99-outputs.tf`: only `outputs` blocks

## Formatting

* Blocks must be spaced by a single empty line
* There must be a single empty line at the end of every file

We also recommend enforcing formatting via `terraform fmt`.

## Rules & configuration

To-do

## To-do

* Define each rule & configuration parameters in the README
* Locals: enforce a single `locals` block
* Setup: rename into `terraform` for consistency

## License

This project is licensed under the [MIT License](LICENSE).
