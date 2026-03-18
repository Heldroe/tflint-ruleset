package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestVariableArgumentsRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid: correct full order",
			files: map[string]string{
				"main.tf": `
variable "example" {
  type        = string
  nullable    = false
  sensitive   = true
  ephemeral   = true
  default     = "foo"
  description = "An example variable"

  validation {
    condition     = length(var.example) > 0
    error_message = "Must not be empty"
  }
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: partial order (type, default, description)",
			files: map[string]string{
				"main.tf": `
variable "example" {
  type        = string
  default     = "foo"
  description = "An example"
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: only description",
			files: map[string]string{
				"main.tf": `
variable "example" {
  description = "Just a description"
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: only type and validation",
			files: map[string]string{
				"main.tf": `
variable "example" {
  type = string

  validation {
    condition     = length(var.example) > 0
    error_message = "Must not be empty"
  }
}
`,
			},
			expected: 0,
		},
		{
			name: "valid: empty variable block",
			files: map[string]string{
				"main.tf": `variable "example" {}`,
			},
			expected: 0,
		},
		{
			name: "valid: unknown arguments are ignored",
			files: map[string]string{
				"main.tf": `
variable "example" {
  type    = string
  default = "foo"
  custom  = "ignored"
}
`,
			},
			expected: 0,
		},
		{
			name: "invalid: description before type",
			files: map[string]string{
				"main.tf": `
variable "example" {
  description = "An example"
  type        = string
}
`,
			},
			expected: 1,
			messages: []string{"'type' must be declared before 'description' in variable blocks (expected order: type, description)"},
		},
		{
			name: "invalid: default before type",
			files: map[string]string{
				"main.tf": `
variable "example" {
  default = "foo"
  type    = string
}
`,
			},
			expected: 1,
			messages: []string{"'type' must be declared before 'default' in variable blocks (expected order: type, default)"},
		},
		{
			name: "invalid: validation before description",
			files: map[string]string{
				"main.tf": `
variable "example" {
  type = string

  validation {
    condition     = length(var.example) > 0
    error_message = "Must not be empty"
  }

  description = "An example"
}
`,
			},
			expected: 1,
			messages: []string{"'description' must be declared before 'validation' in variable blocks (expected order: type, description, validation)"},
		},
		{
			name: "invalid: sensitive after default",
			files: map[string]string{
				"main.tf": `
variable "example" {
  type      = string
  default   = "foo"
  sensitive = true
}
`,
			},
			expected: 1,
			messages: []string{"'sensitive' must be declared before 'default' in variable blocks (expected order: type, sensitive, default)"},
		},
		{
			name: "invalid: multiple ordering violations",
			files: map[string]string{
				"main.tf": `
variable "example" {
  description = "An example"
  default     = "foo"
  type        = string
}
`,
			},
			expected: 2,
			messages: []string{
				"'default' must be declared before 'description' in variable blocks (expected order: type, default, description)",
				"'type' must be declared before 'description' in variable blocks (expected order: type, default, description)",
			},
		},
		{
			name: "multiple variables: one valid one invalid",
			files: map[string]string{
				"main.tf": `
variable "good" {
  type    = string
  default = "ok"
}

variable "bad" {
  default = "not ok"
  type    = string
}
`,
			},
			expected: 1,
			messages: []string{"'type' must be declared before 'default' in variable blocks (expected order: type, default)"},
		},
		{
			name: "non-variable blocks are ignored",
			files: map[string]string{
				"main.tf": `
resource "null_resource" "example" {
  description = "first"
  type        = "second"
}
`,
			},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewVariableArgumentsRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
