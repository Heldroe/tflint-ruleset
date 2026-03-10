package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestFileEndNewlineRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid newline (exactly one at end)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"foo\" {}\n",
			},
			expected: 0,
		},
		{
			name: "missing newline at end",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"foo\" {}",
			},
			expected: 1,
		},
		{
			name: "too many newlines at end (two)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"foo\" {}\n\n",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewFileEndNewlineRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}

func TestBlockSpacingRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid spacing (one blank line between blocks)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 0,
		},
		{
			name: "invalid spacing (no blank line between blocks)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 1,
		},
		{
			name: "invalid spacing (too many blank lines - two)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewBlockSpacingRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}
