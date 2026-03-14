# terraform_style_no_provider_argument

Forbids the use of the `provider` argument within resource and data blocks.

## Details

This rule is **disabled by default**. It is useful if you want to strictly enforce that all resources use the default provider configuration or that provider aliases are handled through other means.

## Why?

Explicit `provider` arguments can make code less portable and harder to reuse in different contexts.

## Examples

### Denied

```hcl
resource "aws_instance" "foo" {
  provider = aws.east
  ami      = "ami-123456"
}
```

### Accepted

```hcl
resource "aws_instance" "foo" {
  ami = "ami-123456"
}
```
