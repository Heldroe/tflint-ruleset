# terraform_style_structure_layout

Enforces consistent layout for lists and maps to prevent "outrageously ugly" code.

## Details

- **Single-line structures** must not have a trailing comma.
- **Multi-line structures** must have the closing bracket (`]` or `}`) on its own line.
- **Multi-line structures** must have their first element on a new line (not on the same line as the opening bracket).
- **Commas** must not be at the start of a line.

## Examples

### Denied

```hcl
# Trailing comma in single line
subnets = ["subnet-1", "subnet-2", ]

# Closing bracket same line as content
availability_zones = [
  "us-east-1a",
  "us-east-1b", ]

# First element same line as opening bracket
public_ips = ["1.1.1.1",
  "2.2.2.2",
]

# Comma at start of line
dns_servers = ["8.8.8.8"
, "8.8.4.4"]
```

### Accepted

```hcl
# Single line list
subnets = ["subnet-1", "subnet-2"]

# Multi-line list
availability_zones = [
  "us-east-1a",
  "us-east-1b",
]

# Nested structures
routes = [
  {
    cidr_block = "0.0.0.0/0"
    gateway_id = "igw-123"
  },
  { cidr_block = "10.0.0.0/16" },
]
```
