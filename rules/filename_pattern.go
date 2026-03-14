package rules

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type FilenamePatternRule struct {
	tflint.DefaultRule
}

func NewFilenamePatternRule() *FilenamePatternRule {
	return &FilenamePatternRule{}
}

func (r *FilenamePatternRule) Name() string {
	return config.RulePrefix + "_filename_pattern"
}

func (r *FilenamePatternRule) Enabled() bool {
	return true
}

func (r *FilenamePatternRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *FilenamePatternRule) Link() string {
	return ruleLink("filename_pattern")
}

func (r *FilenamePatternRule) Check(runner tflint.Runner) error {
	ruleConfig := struct {
		MinIndex int `hclext:"min_index,optional"`
	}{
		MinIndex: 10,
	}

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	// Standard files that are allowed regardless of index
	allowedFiles := map[string]bool{
		config.DefaultVariablesFileName: true, // 00-variables
		config.DefaultTerraformFileName: true, // 01-terraform
		config.DefaultLocalsFileName:    true, // 02-locals
		config.DefaultDataFileName:      true, // 03-data
		config.DefaultOutputsFileName:   true, // 99-outputs
	}

	// Check for custom filenames in other rules
	otherRules := []struct {
		Name            string
		DefaultFilename string
	}{
		{Name: config.RulePrefix + "_variables_file", DefaultFilename: config.DefaultVariablesFileName},
		{Name: config.RulePrefix + "_terraform_file", DefaultFilename: config.DefaultTerraformFileName},
		{Name: config.RulePrefix + "_locals_file", DefaultFilename: config.DefaultLocalsFileName},
		{Name: config.RulePrefix + "_data_file", DefaultFilename: config.DefaultDataFileName},
		{Name: config.RulePrefix + "_outputs_file", DefaultFilename: config.DefaultOutputsFileName},
	}

	for _, otherRule := range otherRules {
		var otherRuleConfig struct {
			Filename string `hclext:"filename,optional"`
		}
		// Initialize with default
		otherRuleConfig.Filename = otherRule.DefaultFilename

		// Decode config
		if err := runner.DecodeRuleConfig(otherRule.Name, &otherRuleConfig); err == nil {
			if otherRuleConfig.Filename != "" {
				allowedFiles[otherRuleConfig.Filename] = true
			}
		}
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`^(\d{2})-[a-z0-9-]+\.tf$`)

	for name := range files {
		baseName := filepath.Base(name)

		// Skip non-tf files if any (though runner usually filters)
		if !strings.HasSuffix(baseName, ".tf") {
			continue
		}

		matches := re.FindStringSubmatch(baseName)
		if matches == nil {
			runner.EmitIssue(
				r,
				"terraform file name must match XX-name.tf",
				hcl.Range{
					Filename: name,
					Start:    hcl.Pos{Line: 1, Column: 1},
					End:      hcl.Pos{Line: 1, Column: 1},
				},
			)
			continue
		}

		// Extract index
		indexStr := matches[1]
		index, _ := strconv.Atoi(indexStr) // regex \d{2} ensures valid int

		// Check if it's a standard file exception
		nameWithoutExt := strings.TrimSuffix(baseName, ".tf")
		if allowedFiles[nameWithoutExt] {
			continue
		}

		// Check min index
		if index < ruleConfig.MinIndex {
			runner.EmitIssue(
				r,
				fmt.Sprintf("custom file index must be >= %d (found %02d)", ruleConfig.MinIndex, index),
				hcl.Range{
					Filename: name,
					Start:    hcl.Pos{Line: 1, Column: 1},
					End:      hcl.Pos{Line: 1, Column: 1},
				},
			)
		}
	}

	return nil
}
