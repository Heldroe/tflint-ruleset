package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleRules(t *testing.T) {
	t.Run("NoEmptyFileRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: file with block",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {}"},
				expected: 0,
			},
			{
				name:     "invalid: empty file",
				files:    map[string]string{"main.tf": ""},
				expected: 1,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewNoEmptyFileRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("TrailingCommaRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: single line list",
				files:    map[string]string{"main.tf": "l = [\"a\", \"b\"]"},
				expected: 0,
			},
			{
				name:     "valid: multi line list with comma",
				files:    map[string]string{"main.tf": "l = [\n\"a\",\n\"b\",\n]"},
				expected: 0,
			},
			{
				name:     "invalid: multi line list missing comma",
				files:    map[string]string{"main.tf": "l = [\n\"a\",\n\"b\"\n]"},
				expected: 1,
			},
			{
				name:     "valid: single line map",
				files:    map[string]string{"main.tf": "m = {a=1, b=2}"},
				expected: 0,
			},
			{
				name:     "valid: multi line map no comma",
				files:    map[string]string{"main.tf": "m = {\na = 1\nb = 2\n}"},
				expected: 0,
			},
			{
				name:     "invalid: multi line map with comma",
				files:    map[string]string{"main.tf": "m = {\na = 1,\nb = 2\n}"},
				expected: 1,
			},
			{
				name:     "invalid: map in function call with comma",
				files:    map[string]string{"main.tf": "x = templatefile(\"path\", {\na = 1,\nb = 2\n})"},
				expected: 1,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewTrailingCommaRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("MapAssignmentRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: use equal",
				files:    map[string]string{"main.tf": "m = {a = 1}"},
				expected: 0,
			},
			{
				name:     "invalid: use colon",
				files:    map[string]string{"main.tf": "m = {a : 1}"},
				expected: 1,
			},
			{
				name:     "invalid: map in function call with colon",
				files:    map[string]string{"main.tf": "x = templatefile(\"path\", {a : 1})"},
				expected: 1,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewMapAssignmentRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("CommentStyleRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: hash with space",
				files:    map[string]string{"main.tf": "# comment"},
				expected: 0,
			},
			{
				name:     "valid: hash only",
				files:    map[string]string{"main.tf": "####\n"},
				expected: 0,
			},
			{
				name:     "valid: five hashes",
				files:    map[string]string{"main.tf": "#####"},
				expected: 0,
			},
			{
				name:     "invalid: slash slash",
				files:    map[string]string{"main.tf": "// comment"},
				expected: 1,
			},
			{
				name:     "invalid: block comment",
				files:    map[string]string{"main.tf": "/* comment */"},
				expected: 1,
			},
			{
				name:     "invalid: hash no space",
				files:    map[string]string{"main.tf": "#comment"},
				expected: 1,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewCommentStyleRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("NoProviderArgumentRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: no provider arg",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {}"},
				expected: 0,
			},
			{
				name:     "invalid: has provider arg",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" { provider = aws }"},
				expected: 1,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewNoProviderArgumentRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("BlockInternalSpacingRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: standard",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n}"},
				expected: 0,
			},
			{
				name:     "invalid: top spacing",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n\n  attr = 1\n}"},
				expected: 1,
			},
			{
				name:     "invalid: bottom spacing",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n\n}"},
				expected: 1,
			},
			{
				name:     "invalid: consecutive blank",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n\n\n  attr2 = 2\n}"},
				expected: 1,
			},
			{
				name:     "valid: comment at top",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  # comment\n  attr = 1\n}"},
				expected: 0,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewBlockInternalSpacingRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("ResourceArgumentsRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: order",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  count = 1\n\n  attr = 1\n\n  lifecycle {}\n\n  depends_on = []\n}"},
				expected: 0,
			},
			{
				name:     "invalid: count not first",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n  count = 1\n}"},
				expected: 1,
			},
			{
				name:     "invalid: no blank after count",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  count = 1\n  attr = 1\n}"},
				expected: 1,
			},
			{
				name:     "invalid: depends_on not last",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  depends_on = []\n  attr = 1\n}"},
				expected: 1,
			},
			{
				name:     "invalid: no blank above depends_on",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n  depends_on = []\n}"},
				expected: 1,
			},
			{
				name:     "invalid: lifecycle after depends_on",
				files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  depends_on = []\n  lifecycle {}\n}"},
				expected: 2, // depends_on not last, lifecycle not before depends_on
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewResourceArgumentsRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})

	t.Run("StructureLayoutRule", func(t *testing.T) {
		tests := []struct {
			name     string
			files    map[string]string
			expected int
		}{
			{
				name:     "valid: single line list",
				files:    map[string]string{"main.tf": "l = [1, 2]"},
				expected: 0,
			},
			{
				name:     "invalid: single line list trailing comma",
				files:    map[string]string{"main.tf": "l = [1, 2, ]"},
				expected: 1,
			},
			{
				name:     "invalid: multi line list closing bracket same line",
				files:    map[string]string{"main.tf": "l = [\n  1,\n  2, ]"},
				expected: 1,
			},
			{
				name:     "invalid: comma at start of line",
				files:    map[string]string{"main.tf": "l = [1\n, 2]"},
				expected: 2,
			},
			{
				name:     "valid: multi line list correct",
				files:    map[string]string{"main.tf": "l = [\n  1,\n  2,\n]"},
				expected: 0,
			},
			{
				name:     "valid: list of single line maps",
				files:    map[string]string{"main.tf": "l = [\n  { foo = \"bar\" },\n  { foo = \"baz\" },\n]"},
				expected: 0,
			},
			{
				name:     "invalid: multi line map closing bracket same line",
				files:    map[string]string{"main.tf": "l = [\n  {\n    foo = \"bar\" },\n]"},
				expected: 1,
			},
			{
				name:     "valid: multi line map correct",
				files:    map[string]string{"main.tf": "l = [\n  {\n    foo = \"bar\"\n  },\n]"},
				expected: 0,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				runner := helper.TestRunner(t, tc.files)
				rule := rules.NewStructureLayoutRule()
				if err := rule.Check(runner); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if len(runner.Issues) != tc.expected {
					t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
				}
			})
		}
	})
}
