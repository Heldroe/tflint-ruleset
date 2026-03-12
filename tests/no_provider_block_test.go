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
	}{
		{
			name: "valid: no provider block",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"foo\" {}",
			},
			expected: 0,
		},
		{
			name: "invalid: provider block in main.tf",
			files: map[string]string{
				"main.tf": "provider \"aws\" {}",
			},
			expected: 1,
		},
		{
			name: "invalid: provider block in provider.tf",
			files: map[string]string{
				"provider.tf": "provider \"aws\" {}",
			},
			expected: 1,
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewNoProviderBlockRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}
