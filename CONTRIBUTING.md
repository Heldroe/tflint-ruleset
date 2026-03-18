# Contributing to TFLint Ruleset for Terraform Style

Thank you for your interest in contributing to this project!

## Development

### Prerequisites

- [Go](https://golang.org/doc/install) (see `go.mod` for the required version)
- [TFLint](https://github.com/terraform-linters/tflint)

### Building and Testing

You can use the provided `Makefile` for common tasks:

```bash
# Build the plugin binary
make build

# Run all tests
make test

# Format the code
make fmt

# Install the plugin locally (~/.tflint.d/plugins)
make install
```

## Release Process

This project uses [GoReleaser](https://goreleaser.com/) and GitHub Actions for releases. To create a new release, follow these steps:

1.  **Update the Version**: Update the `version` variable in `ruleset.go` to the new version (e.g., `0.2.0`). Also update the version references in `README.md` (installation and recommended configuration sections) to match.
2.  **Commit and Push**: Commit this change to the `main` branch.
3.  **Tag the Release**: Create a new git tag that matches the version in `ruleset.go`, prefixed with `v`.
    ```bash
    git tag v0.2.0
    git push origin v0.2.0
    ```
4.  **Verification**: The GitHub Actions release workflow will automatically trigger. It includes a verification step that ensures the git tag matches the version string in `ruleset.go`. If they do not match, the release will fail.
5.  **Build Injection**: During the build process, the version is also injected into the binary via `ldflags` to ensure consistency.

Once the CI workflow completes, the new release will be available on the GitHub Releases page with automatically generated changelogs and signed binaries.
