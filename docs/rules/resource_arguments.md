# terraform_style_resource_arguments

Enforces ordering and spacing of standard arguments and blocks within resources and modules.

## Details

- `count` and `for_each` must be the first things declared in a block.
- There must be an empty blank line after `count` or `for_each`.
- In `module` blocks, `source` must be the first argument (but below `count` and `for_each` if present).
- `lifecycle` must be the last thing declared (but before `depends_on` if present).
- `depends_on` must be the last thing declared.
- There must be an empty blank line above `depends_on`.

## Examples

### Denied

```hcl
resource "aws_instance" "foo" {
  ami = "ami-123456"
  count = 1 # count not first
  depends_on = [aws_vpc.main] # missing blank line above
}
```

### Accepted

```hcl
resource "aws_instance" "foo" {
  count = 1

  ami           = "ami-123456"
  instance_type = "t2.micro"

  lifecycle {
    ignore_changes = [ami]
  }

  depends_on = [aws_vpc.main]
}
```
