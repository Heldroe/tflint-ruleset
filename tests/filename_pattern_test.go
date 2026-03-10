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
			name: "valid filename",
			files: map[string]string{
				"00-main.tf": "",
			},
			expected: 0,
		},
		{
			name: "invalid filename - no prefix (e.g. main.tf)",
			files: map[string]string{
				"main.tf": "",
			},
			expected: 1,
		},
		{
			name: "invalid filename - uppercase (e.g. 01-SETUP.tf)",
			files: map[string]string{
				"01-SETUP.tf": "",
			},
			expected: 1,
		},
		{
			name: "invalid filename - underscore (e.g. 05_locals.tf)",
			files: map[string]string{
				"05_locals.tf": "",
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
