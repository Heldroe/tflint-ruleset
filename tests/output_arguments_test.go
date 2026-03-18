package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestOutputArgumentsRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid: correct full order",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  description = "An example output"
  sensitive   = true
  ephemeral   = true
  value       = "foo"

  precondition {
    condition     = length(var.example) > 0
    error_message = "Must not be empty"
  }

  depends_on = [null_resource.example]
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: partial order (description, value)",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  description = "An example"
  value       = "foo"
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: only value",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  value = "foo"
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: empty output block",
			files: map[string]string{
				"99-outputs.tf": `output "example" {}`,
			},
			expected: 0,
		},
		{
			name: "valid: unknown arguments are ignored",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  description = "An example"
  value       = "foo"
  custom      = "ignored"
}
`,
			},
			expected: 0,
		},
		{
			name: "invalid: value before description",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  value       = "foo"
  description = "An example"
}
`,
			},
			expected: 1,
			messages: []string{"'description' must be declared before 'value' in output blocks (expected order: description, value)"},
		},
		{
			name: "invalid: depends_on before value",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  depends_on = [null_resource.example]
  value      = "foo"
}
`,
			},
			expected: 1,
			messages: []string{"'value' must be declared before 'depends_on' in output blocks (expected order: value, depends_on)"},
		},
		{
			name: "invalid: precondition before value",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  description = "An example"

  precondition {
    condition     = true
    error_message = "fail"
  }

  value = "foo"
}
`,
			},
			expected: 1,
			messages: []string{"'value' must be declared before 'precondition' in output blocks (expected order: description, value, precondition)"},
		},
		{
			name: "invalid: multiple violations",
			files: map[string]string{
				"99-outputs.tf": `
output "example" {
  sensitive   = true
  depends_on  = [null_resource.example]
  description = "An example"
  value       = "foo"
}
`,
			},
			expected: 2,
			messages: []string{
				"'description' must be declared before 'depends_on' in output blocks (expected order: description, sensitive, value, depends_on)",
				"'value' must be declared before 'depends_on' in output blocks (expected order: description, sensitive, value, depends_on)",
			},
		},
		{
			name: "non-output blocks are ignored",
			files: map[string]string{
				"main.tf": `
resource "null_resource" "example" {
  value       = "first"
  description = "second"
}
`,
			},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewOutputArgumentsRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
