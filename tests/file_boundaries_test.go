package rules_test

import (
	"testing"

	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestVariablesFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid variables location",
			files: map[string]string{
				"00-variables.tf": `variable "foo" {}`,
			},
			expected: 0,
		},
		{
			name: "valid multiple variables in same file",
			files: map[string]string{
				"00-variables.tf": `variable "foo" {}
variable "bar" {}`,
			},
			expected: 0,
		},
		{
			name: "invalid variables location (in main.tf)",
			files: map[string]string{
				"main.tf": `variable "foo" {}`,
			},
			expected: 1,
			messages: []string{"variable blocks must be defined in 00-variables.tf"},
		},
		{
			name: "unauthorized block in variables file (resource in 00-variables.tf)",
			files: map[string]string{
				"00-variables.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only variable blocks are allowed in 00-variables.tf; found resource"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewVariablesFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestTerraformBlockFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid terraform block location",
			files: map[string]string{
				"01-terraform.tf": `terraform { required_version = ">= 1.0" }`,
			},
			expected: 0,
		},
		{
			name: "invalid terraform block location (in main.tf)",
			files: map[string]string{
				"main.tf": `terraform { required_version = ">= 1.0" }`,
			},
			expected: 1,
			messages: []string{"terraform blocks must be defined in 01-terraform.tf"},
		},
		{
			name: "multiple terraform blocks in terraform file",
			files: map[string]string{
				"01-terraform.tf": `
terraform {
  required_version = ">= 1.0"
}
terraform {
  required_providers {
    null = {}
  }
}
`,
			},
			expected: 1,
			messages: []string{"only 1 terraform block(s) allowed in 01-terraform.tf; found multiple"},
		},
		{
			name: "unauthorized block in terraform file",
			files: map[string]string{
				"01-terraform.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only terraform blocks are allowed in 01-terraform.tf; found resource"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewTerraformBlockFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestOutputsFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid outputs location",
			files: map[string]string{
				"99-outputs.tf": `output "foo" { value = "bar" }`,
			},
			expected: 0,
		},
		{
			name: "invalid outputs location (in main.tf)",
			files: map[string]string{
				"main.tf": `output "foo" { value = "bar" }`,
			},
			expected: 1,
			messages: []string{"output blocks must be defined in 99-outputs.tf"},
		},
		{
			name: "unauthorized block in outputs file",
			files: map[string]string{
				"99-outputs.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only output blocks are allowed in 99-outputs.tf; found resource"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewOutputsFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestLocalsFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid locals location",
			files: map[string]string{
				"02-locals.tf": `locals { foo = "bar" }`,
			},
			expected: 0,
		},
		{
			name: "invalid locals location (in main.tf)",
			files: map[string]string{
				"main.tf": `locals { foo = "bar" }`,
			},
			expected: 1,
			messages: []string{"locals blocks must be defined in 02-locals.tf"},
		},
		{
			name: "multiple locals blocks in locals file",
			files: map[string]string{
				"02-locals.tf": `
locals { foo = "bar" }
locals { baz = "qux" }
`,
			},
			expected: 1,
			messages: []string{"only 1 locals block(s) allowed in 02-locals.tf; found multiple"},
		},
		{
			name: "unauthorized block in locals file",
			files: map[string]string{
				"02-locals.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only locals blocks are allowed in 02-locals.tf; found resource"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewLocalsFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}

func TestDataFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid data location",
			files: map[string]string{
				"03-data.tf": `data "null_data_source" "foo" {}`,
			},
			expected: 0,
		},
		{
			name: "invalid data location (in main.tf)",
			files: map[string]string{
				"main.tf": `data "null_data_source" "foo" {}`,
			},
			expected: 1,
			messages: []string{"data blocks must be defined in 03-data.tf"},
		},
		{
			name: "unauthorized block in data file",
			files: map[string]string{
				"03-data.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only data blocks are allowed in 03-data.tf; found resource"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewDataFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
