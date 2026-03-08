package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestVariablesFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid variables location",
			files: map[string]string{
				"00-variables.tf": "variable \"foo\" {}",
			},
			expected: 0,
		},
		{
			name: "invalid variables location (in main.tf)",
			files: map[string]string{
				"main.tf": "variable \"foo\" {}",
			},
			expected: 1,
		},
		{
			name: "unauthorized block in variables file (resource in 00-variables.tf)",
			files: map[string]string{
				"00-variables.tf": "resource \"null_resource\" \"foo\" {}",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := NewVariablesFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}

func TestTerraformBlockFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid terraform block location",
			files: map[string]string{
				"01-setup.tf": "terraform { required_version = \">= 1.0\" }",
			},
			expected: 0,
		},
		{
			name: "invalid terraform block location (in main.tf)",
			files: map[string]string{
				"main.tf": "terraform { required_version = \">= 1.0\" }",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := NewTerraformBlockFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}

func TestOutputsFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid outputs location",
			files: map[string]string{
				"99-outputs.tf": "output \"foo\" { value = \"bar\" }",
			},
			expected: 0,
		},
		{
			name: "invalid outputs location (in main.tf)",
			files: map[string]string{
				"main.tf": "output \"foo\" { value = \"bar\" }",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := NewOutputsFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}

func TestLocalsFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid locals location",
			files: map[string]string{
				"05-locals.tf": "locals { foo = \"bar\" }",
			},
			expected: 0,
		},
		{
			name: "invalid locals location (in main.tf)",
			files: map[string]string{
				"main.tf": "locals { foo = \"bar\" }",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := NewLocalsFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}

func TestDataFileRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
	}{
		{
			name: "valid data location",
			files: map[string]string{
				"10-data.tf": "data \"null_data_source\" \"foo\" {}",
			},
			expected: 0,
		},
		{
			name: "invalid data location (in main.tf)",
			files: map[string]string{
				"main.tf": "data \"null_data_source\" \"foo\" {}",
			},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			rule := NewDataFileRule()

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(runner.Issues) != tc.expected {
				t.Errorf("%s: expected %d issues, got %d", tc.name, tc.expected, len(runner.Issues))
			}
		})
	}
}
