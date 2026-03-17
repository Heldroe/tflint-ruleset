package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func assertIssues(t *testing.T, runner *helper.Runner, expected int, messages []string) {
	t.Helper()
	if len(runner.Issues) != expected {
		t.Errorf("expected %d issues, got %d", expected, len(runner.Issues))
		for i, issue := range runner.Issues {
			t.Logf("  issue[%d]: %s", i, issue.Message)
		}
		return
	}
	for i, msg := range messages {
		if i >= len(runner.Issues) {
			break
		}
		if runner.Issues[i].Message != msg {
			t.Errorf("issue[%d]: expected message %q, got %q", i, msg, runner.Issues[i].Message)
		}
	}
}

func TestNoEmptyFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: file with block",
			files:    map[string]string{"main.tf": `resource "x" "y" {}`},
			expected: 0,
		},
		{
			name:     "valid: file with attribute",
			files:    map[string]string{"main.tf": `foo = "bar"`},
			expected: 0,
		},
		{
			name:     "invalid: empty file",
			files:    map[string]string{"main.tf": ""},
			expected: 1,
			messages: []string{"file cannot be empty"},
		},
		{
			name:     "invalid: whitespace only file",
			files:    map[string]string{"main.tf": "   \n\n  "},
			expected: 1,
			messages: []string{"file cannot be empty"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewNoEmptyFileRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestTrailingCommaRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: single line list",
			files:    map[string]string{"main.tf": `l = ["a", "b"]`},
			expected: 0,
		},
		{
			name:     "valid: multi line list with comma",
			files:    map[string]string{"main.tf": "l = [\n\"a\",\n\"b\",\n]"},
			expected: 0,
		},
		{
			name:     "valid: empty list",
			files:    map[string]string{"main.tf": "l = []"},
			expected: 0,
		},
		{
			name:     "invalid: multi line list missing comma",
			files:    map[string]string{"main.tf": "l = [\n\"a\",\n\"b\"\n]"},
			expected: 1,
			messages: []string{"lists defined in multiple lines must have a trailing comma on the last line"},
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
			messages: []string{"maps defined in multiple lines must not have any trailing comma on any lines"},
		},
		{
			name:     "invalid: map in function call with comma",
			files:    map[string]string{"main.tf": "x = templatefile(\"path\", {\na = 1,\nb = 2\n})"},
			expected: 1,
			messages: []string{"maps defined in multiple lines must not have any trailing comma on any lines"},
		},
		{
			name: "valid: single-element multi-line list without comma (with exclude_single_element=true)",
			files: map[string]string{
				"main.tf": "l = [\n  1\n]",
				".tflint.hcl": `
rule "terraform_style_trailing_comma" {
  enabled = true
  exclude_single_element = true
}
`,
			},
			expected: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewTrailingCommaRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestMapAssignmentRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: use equal",
			files:    map[string]string{"main.tf": "m = {a = 1}"},
			expected: 0,
		},
		{
			name:     "valid: nested map with equal",
			files:    map[string]string{"main.tf": "m = {\na = {\n  b = 1\n}\n}"},
			expected: 0,
		},
		{
			name:     "invalid: use colon",
			files:    map[string]string{"main.tf": "m = {a : 1}"},
			expected: 1,
			messages: []string{"maps must use the equal '=' sign to assign keys to values"},
		},
		{
			name:     "invalid: map in function call with colon",
			files:    map[string]string{"main.tf": `x = templatefile("path", {a : 1})`},
			expected: 1,
			messages: []string{"maps must use the equal '=' sign to assign keys to values"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewMapAssignmentRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestCommentStyleRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
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
			name:     "valid: empty comment",
			files:    map[string]string{"main.tf": "#\n"},
			expected: 0,
		},
		{
			name:     "invalid: slash slash",
			files:    map[string]string{"main.tf": "// comment"},
			expected: 1,
			messages: []string{"comments are only allowed using '#' (number sign)"},
		},
		{
			name:     "invalid: block comment",
			files:    map[string]string{"main.tf": "/* comment */"},
			expected: 1,
			messages: []string{"comments are only allowed using '#' (number sign)"},
		},
		{
			name:     "invalid: hash no space",
			files:    map[string]string{"main.tf": "#comment"},
			expected: 1,
			messages: []string{"there must be a single space between the '#' and the beginning of the comment"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewCommentStyleRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestNoProviderArgumentRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: no provider arg",
			files:    map[string]string{"main.tf": `resource "x" "y" {}`},
			expected: 0,
		},
		{
			name:     "invalid: has provider arg",
			files:    map[string]string{"main.tf": `resource "x" "y" { provider = aws }`},
			expected: 1,
			messages: []string{`the "provider" argument is not allowed`},
		},
		{
			name: "invalid: nested provider arg",
			files: map[string]string{"main.tf": `
resource "x" "y" {
  provisioner "local-exec" {
    provider = aws
  }
}
`},
			expected: 1,
			messages: []string{`the "provider" argument is not allowed`},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewNoProviderArgumentRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestBlockInternalSpacingRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: standard",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n}"},
			expected: 0,
		},
		{
			name:     "valid: empty block",
			files:    map[string]string{"main.tf": `resource "x" "y" {}`},
			expected: 0,
		},
		{
			name:     "valid: comment at top",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  # comment\n  attr = 1\n}"},
			expected: 0,
		},
		{
			name:     "valid: single blank line between items",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n\n  attr2 = 2\n}"},
			expected: 0,
		},
		{
			name:     "invalid: top spacing",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n\n  attr = 1\n}"},
			expected: 1,
			messages: []string{"there can't be a spacing empty line at the very top of the block"},
		},
		{
			name:     "invalid: bottom spacing",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n\n}"},
			expected: 1,
			messages: []string{"there can't be a spacing empty line at the very bottom of the block"},
		},
		{
			name:     "invalid: consecutive blank",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n\n\n  attr2 = 2\n}"},
			expected: 1,
			messages: []string{"there cannot be more than 1 empty line of spacing in a row within blocks"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewBlockInternalSpacingRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestResourceArgumentsRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: correct order",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  count = 1\n\n  attr = 1\n\n  lifecycle {}\n\n  depends_on = []\n}"},
			expected: 0,
		},
		{
			name:     "valid: for_each first",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  for_each = {}\n\n  attr = 1\n}"},
			expected: 0,
		},
		{
			name:     "invalid: count not first",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n  count = 1\n}"},
			expected: 1,
			messages: []string{"'count' and 'for_each' arguments must be the first thing declared in the block"},
		},
		{
			name:     "invalid: no blank after count",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  count = 1\n  attr = 1\n}"},
			expected: 1,
			messages: []string{"there must be an empty blank line after count"},
		},
		{
			name:     "invalid: depends_on not last",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  depends_on = []\n  attr = 1\n}"},
			expected: 1,
			messages: []string{"'depends_on' must be the last thing declared in the block"},
		},
		{
			name:     "invalid: no blank above depends_on",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  attr = 1\n  depends_on = []\n}"},
			expected: 1,
			messages: []string{"there must be an empty blank line above depends_on"},
		},
		{
			name:     "invalid: lifecycle after depends_on",
			files:    map[string]string{"main.tf": "resource \"x\" \"y\" {\n  depends_on = []\n  lifecycle {}\n}"},
			expected: 2,
		},
		{
			name:     "valid: module with source first",
			files:    map[string]string{"main.tf": "module \"x\" {\n  source = \"./mod\"\n\n  attr = 1\n}"},
			expected: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewResourceArgumentsRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestStructureLayoutRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name:     "valid: single line list",
			files:    map[string]string{"main.tf": "l = [1, 2]"},
			expected: 0,
		},
		{
			name:     "valid: empty list",
			files:    map[string]string{"main.tf": "l = []"},
			expected: 0,
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
			name:     "valid: multi line map correct",
			files:    map[string]string{"main.tf": "l = [\n  {\n    foo = \"bar\"\n  },\n]"},
			expected: 0,
		},
		{
			name:     "invalid: single line list trailing comma",
			files:    map[string]string{"main.tf": "l = [1, 2, ]"},
			expected: 1,
			messages: []string{"single-line structures must not have a trailing comma"},
		},
		{
			name:     "invalid: multi line list closing bracket same line",
			files:    map[string]string{"main.tf": "l = [\n  1,\n  2, ]"},
			expected: 1,
			messages: []string{"multi-line structures must have the closing bracket on its own line"},
		},
		{
			name:     "invalid: comma at start of line",
			files:    map[string]string{"main.tf": "l = [1\n, 2]"},
			expected: 3,
		},
		{
			name:     "invalid: multi line map closing bracket same line",
			files:    map[string]string{"main.tf": "l = [\n  {\n    foo = \"bar\" },\n]"},
			expected: 1,
			messages: []string{"multi-line structures must have the closing bracket on its own line"},
		},
		{
			name:     "invalid: multi line list first element same line",
			files:    map[string]string{"main.tf": "l = [1,\n  2,\n]"},
			expected: 1,
			messages: []string{"for multi-line structures, the first element must be on a new line"},
		},
		{
			name:     "valid: single line map",
			files:    map[string]string{"main.tf": "m = {a = 1, b = 2}"},
			expected: 0,
		},
		{
			name:     "invalid: single line map trailing comma",
			files:    map[string]string{"main.tf": "m = {a = 1, }"},
			expected: 1,
			messages: []string{"single-line structures must not have a trailing comma"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewStructureLayoutRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
