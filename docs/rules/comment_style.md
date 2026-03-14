# terraform_style_comment_style

Enforces a consistent comment style in Terraform files.

## Details

- Only single-line comments starting with `#` are allowed.
- `//` and `/* ... */` comments are prohibited.
- There must be a single space between the `#` and the beginning of the comment text.
- Exception: Comments consisting only of `#` signs (e.g., `#######`) are allowed and don't require a space.

## Why?

Consistent commenting improves code readability and follows the recommended Terraform style.

## Examples

### Denied

```hcl
// Denied comment style
/* Block comment */
#comment without space
```

### Accepted

```hcl
# Accepted comment style
#######
# A separator above
```
