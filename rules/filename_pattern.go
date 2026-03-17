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

	allowedFiles := resolveSpecialFiles(runner)

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
