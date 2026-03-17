package rules

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
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

	allowedFiles := map[string]bool{
		config.DefaultVariablesFileName: true,
		config.DefaultTerraformFileName: true,
		config.DefaultLocalsFileName:    true,
		config.DefaultDataFileName:      true,
		config.DefaultOutputsFileName:   true,
	}

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
		otherRuleConfig.Filename = otherRule.DefaultFilename

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

	for name := range files {
		baseName := filepath.Base(name)

		if !strings.HasSuffix(baseName, ".tf") {
			continue
		}

		matches := filenamePattern.FindStringSubmatch(baseName)
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

		index, _ := strconv.Atoi(matches[1])

		nameWithoutExt := strings.TrimSuffix(baseName, ".tf")
		if allowedFiles[nameWithoutExt] {
			continue
		}

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
