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
		messages []string
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
			messages: []string{"file must end with exactly one newline"},
		},
		{
			name: "too many newlines at end (two)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"foo\" {}\n\n",
			},
			expected: 1,
			messages: []string{"file must end with exactly one newline"},
		},
		{
			name: "too many newlines at end (three)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"foo\" {}\n\n\n",
			},
			expected: 1,
			messages: []string{"file must end with exactly one newline"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewFileEndNewlineRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestBlockSpacingRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid spacing (one blank line between blocks)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 0,
		},
		{
			name: "valid single block (no spacing needed)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n",
			},
			expected: 0,
		},
		{
			name: "invalid spacing (no blank line between blocks)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 1,
			messages: []string{"blocks must be separated by exactly one blank line"},
		},
		{
			name: "invalid spacing (too many blank lines - two)",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 1,
			messages: []string{"blocks must be separated by exactly one blank line"},
		},
		{
			name: "valid spacing with comment between blocks",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\n# comment\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 0,
		},
		{
			name: "valid spacing with multiple comment lines between blocks",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\n# line 1\n# line 2\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 0,
		},
		{
			name: "invalid spacing with comment but too many blank lines",
			files: map[string]string{
				"main.tf": "resource \"null_resource\" \"a\" {}\n\n\n# comment\n\nresource \"null_resource\" \"b\" {}\n",
			},
			expected: 1,
			messages: []string{"blocks must be separated by exactly one blank line"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewBlockSpacingRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
