package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestNoProviderBlockRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid: no provider block",
			files: map[string]string{
				"main.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 0,
		},
		{
			name: "invalid: provider block in main.tf",
			files: map[string]string{
				"main.tf": `provider "aws" {}`,
			},
			expected: 1,
			messages: []string{"provider blocks are not allowed"},
		},
		{
			name: "invalid: provider block in provider.tf",
			files: map[string]string{
				"provider.tf": `provider "aws" {}`,
			},
			expected: 1,
			messages: []string{"provider blocks are not allowed"},
		},
		{
			name: "invalid: multiple provider blocks",
			files: map[string]string{
				"main.tf": `
provider "aws" {
  region = "us-east-1"
}
provider "google" {}
`,
			},
			expected: 2,
			messages: []string{
				"provider blocks are not allowed",
				"provider blocks are not allowed",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewNoProviderBlockRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
