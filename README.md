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

## License

This project is licensed under the [MIT License](LICENSE).
