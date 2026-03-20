# terraform_style_filename_pattern

Enforces that all Terraform files follow a consistent naming convention to ensure a deterministic ordering in file browsers.

## Details

All `.tf` files are expected to match the regex `^(\d{2})-[a-z0-9-]+\.tf$`. This means they must start with two digits, followed by a hyphen, and then a lowercase alphanumeric name.

## Configuration

| Name | Description | Type | Default |
| --- | --- | --- | --- |
| `min_index` | The minimum allowed numeric prefix for custom (non-standard) resource files. | int | `10` |

## Examples

### Denied

```text
main.tf              # Missing numeric prefix
variables.tf         # Missing numeric prefix
0-init.tf            # Prefix must be two digits
10_resource.tf       # Must use hyphen instead of underscore
10-Resource.tf       # Must be lowercase
```

### Accepted

```text
00-variables.tf
01-terraform.tf
10-main.tf
20-network.tf
99-outputs.tf
```
