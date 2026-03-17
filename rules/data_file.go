package rules

import (
	"fmt"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DataFileRule struct {
	tflint.DefaultRule
}

func NewDataFileRule() *DataFileRule {
	return &DataFileRule{}
}

func (r *DataFileRule) Name() string {
	return config.RulePrefix + "_data_file"
}

func (r *DataFileRule) Enabled() bool {
	return true
}

func (r *DataFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *DataFileRule) Link() string {
	return ruleLink("data_file")
}

func (r *DataFileRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		Filename      string   `hclext:"filename,optional"`
		AllowedBlocks []string `hclext:"allowed_blocks,optional"`
	}

	ruleConfig.Filename = config.DefaultDataFileName
	ruleConfig.AllowedBlocks = []string{"data"}

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	expected := fmt.Sprintf("%s.tf", ruleConfig.Filename)
	return enforceFileAllowedBlocks(runner, r, expected, ruleConfig.AllowedBlocks, nil)
}
