# terraform_style_block_internal_spacing

Enforces consistent spacing inside resource, module, and other blocks.

## Details

- There cannot be more than 1 consecutive empty line of spacing within a block.
- There must not be an empty line at the very top of a block (immediately after `{`). Comments are allowed.
- There must not be an empty line at the very bottom of a block (immediately before `}`).

## Examples

### Denied

```hcl
resource "aws_instance" "foo" {

  ami = "ami-123456" # empty line at top


  instance_type = "t2.micro" # 2 empty lines above

} # empty line at bottom
```

### Accepted

```hcl
resource "aws_instance" "foo" {
  # Comments are fine at the top
  ami           = "ami-123456"
  instance_type = "t2.micro"
}
```
