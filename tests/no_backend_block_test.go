package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestNoBackendBlockRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid: no backend block",
			files: map[string]string{
				"01-terraform.tf": `
terraform {
  required_version = ">= 1.0"
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: non-terraform block ignored",
			files: map[string]string{
				"main.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 0,
		},
		{
			name: "invalid: backend block in terraform block",
			files: map[string]string{
				"01-terraform.tf": `
terraform {
  backend "s3" {
    bucket = "my-bucket"
  }
}
`,
			},
			expected: 1,
			messages: []string{"backend configuration is not allowed in the terraform block"},
		},
		{
			name: "invalid: multiple backend blocks",
			files: map[string]string{
				"01-terraform.tf": `
terraform {
  backend "s3" {}
  backend "azurerm" {}
}
`,
			},
			expected: 2,
			messages: []string{
				"backend configuration is not allowed in the terraform block",
				"backend configuration is not allowed in the terraform block",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewNoBackendBlockRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
