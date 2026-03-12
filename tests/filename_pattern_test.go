package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestFilenamePatternRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid filename (index >= 10)",
			files: map[string]string{
				"10-main.tf": "",
			},
			expected: 0,
		},
		{
			name: "valid standard file (00-variables.tf)",
			files: map[string]string{
				"00-variables.tf": "",
			},
			expected: 0,
		},
		{
			name: "valid standard file (01-terraform.tf)",
			files: map[string]string{
				"01-terraform.tf": "",
			},
			expected: 0,
		},
		{
			name: "valid standard file (02-locals.tf)",
			files: map[string]string{
				"02-locals.tf": "",
			},
			expected: 0,
		},
		{
			name: "valid standard file (03-data.tf)",
			files: map[string]string{
				"03-data.tf": "",
			},
			expected: 0,
		},
		{
			name: "valid standard file (99-outputs.tf)",
			files: map[string]string{
				"99-outputs.tf": "",
			},
			expected: 0,
		},
		{
			name: "invalid filename - too low index (e.g. 05-foo.tf)",
			files: map[string]string{
				"05-foo.tf": "",
			},
			expected: 1,
		},
		{
			name: "invalid filename - too low index (e.g. 00-custom.tf)",
			files: map[string]string{
				"00-custom.tf": "",
			},
			expected: 1,
		},
		{
			name: "invalid filename - no prefix (e.g. main.tf)",
			files: map[string]string{
				"main.tf": "",
			},
			expected: 1,
		},
		{
			name: "invalid filename - uppercase (e.g. 10-TERRAFORM.tf)",
			files: map[string]string{
				"10-TERRAFORM.tf": "",
			},
			expected: 1,
		},
		{
			name: "invalid filename - underscore (e.g. 02_locals.tf)",
			files: map[string]string{
				"02_locals.tf": "",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewFilenamePatternRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			issues := runner.Issues
			if len(issues) != tc.expected {
				t.Errorf("expected %d issues, got %d", tc.expected, len(issues))
			}
		})
	}
}
