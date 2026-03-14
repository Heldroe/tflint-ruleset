# terraform_style_block_spacing

Enforces that top-level blocks in a file are separated by exactly one blank line.

## Examples

### Denied

```hcl
# No blank line between blocks
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id
}

# Too many blank lines (more than one)
resource "aws_vpc" "other" {
}


resource "aws_subnet" "private" {
}
```

### Accepted

```hcl
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id
}
```
