# terraform_style_file_end_newline

Enforces that every Terraform file ends with exactly one newline character.

## Why?

Ensuring that every file ends with a newline is a common practice that helps with various command-line tools and version control systems (e.g., preventing `\ No newline at end of file` in git diffs).

## Examples

### Denied

End of file without newline
```hcl
resource "x" "y" {}
```

End of file with two newlines
```hcl
resource "x" "y" {}


```

### Accepted

End of file with exactly one newline
```hcl
resource "x" "y" {}

```
