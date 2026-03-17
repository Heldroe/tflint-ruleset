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
			name: "valid check block in variables file",
			files: map[string]string{
				"00-variables.tf": `check "health" {
  data "http" "example" { url = "https://example.com" }
}`,
			},
			expected: 0,
		},
		{
			name: "unauthorized block in variables file (resource in 00-variables.tf)",
			files: map[string]string{
				"00-variables.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only variable, check blocks are allowed in 00-variables.tf; found resource"},
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

func TestResourceFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		messages []string
	}{
		{
			name: "valid resource block in resource file",
			files: map[string]string{
				"10-main.tf": `resource "null_resource" "foo" {}`,
			},
			expected: 0,
		},
		{
			name: "valid module block in resource file",
			files: map[string]string{
				"10-main.tf": `module "foo" { source = "./mod" }`,
			},
			expected: 0,
		},
		{
			name: "variable block not allowed in resource file",
			files: map[string]string{
				"10-main.tf": `variable "foo" {}`,
			},
			expected: 1,
			messages: []string{"only check, module, moved, removed, resource blocks are allowed in resource files; found variable in 10-main.tf"},
		},
		{
			name: "output block not allowed in resource file",
			files: map[string]string{
				"10-main.tf": `output "foo" { value = "bar" }`,
			},
			expected: 1,
			messages: []string{"only check, module, moved, removed, resource blocks are allowed in resource files; found output in 10-main.tf"},
		},
		{
			name: "data block not allowed in resource file",
			files: map[string]string{
				"10-main.tf": `data "null_data_source" "foo" {}`,
			},
			expected: 1,
			messages: []string{"only check, module, moved, removed, resource blocks are allowed in resource files; found data in 10-main.tf"},
		},
		{
			name: "locals block not allowed in resource file",
			files: map[string]string{
				"10-main.tf": `locals { foo = "bar" }`,
			},
			expected: 1,
			messages: []string{"only check, module, moved, removed, resource blocks are allowed in resource files; found locals in 10-main.tf"},
		},
		{
			name: "provider block not allowed in resource file",
			files: map[string]string{
				"10-main.tf": `provider "aws" { region = "us-east-1" }`,
			},
			expected: 1,
			messages: []string{"only check, module, moved, removed, resource blocks are allowed in resource files; found provider in 10-main.tf"},
		},
		{
			name: "terraform block not allowed in resource file",
			files: map[string]string{
				"10-main.tf": `terraform { required_version = ">= 1.0" }`,
			},
			expected: 1,
			messages: []string{"only check, module, moved, removed, resource blocks are allowed in resource files; found terraform in 10-main.tf"},
		},
		{
			name: "excluded files are skipped",
			files: map[string]string{
				"00-variables.tf": `variable "foo" {}`,
				"02-locals.tf":    `locals { foo = "bar" }`,
				"03-data.tf":      `data "null_data_source" "foo" {}`,
				"99-outputs.tf":   `output "foo" { value = "bar" }`,
			},
			expected: 0,
		},
		{
			name: "multiple disallowed blocks in resource file",
			files: map[string]string{
				"10-main.tf": `
variable "foo" {}
output "bar" { value = "baz" }
`,
			},
			expected: 2,
			messages: []string{
				"only check, module, moved, removed, resource blocks are allowed in resource files; found variable in 10-main.tf",
				"only check, module, moved, removed, resource blocks are allowed in resource files; found output in 10-main.tf",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := rules.NewResourceFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertIssues(t, runner, tc.expected, tc.messages)
		})
	}
}
