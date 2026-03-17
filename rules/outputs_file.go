package rules

import (
	"fmt"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type OutputsFileRule struct {
	tflint.DefaultRule
}

func NewOutputsFileRule() *OutputsFileRule {
	return &OutputsFileRule{}
}

func (r *OutputsFileRule) Name() string {
	return config.RulePrefix + "_outputs_file"
}

func (r *OutputsFileRule) Enabled() bool {
	return true
}

func (r *OutputsFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *OutputsFileRule) Link() string {
	return ruleLink("outputs_file")
}

func (r *OutputsFileRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		Filename string `hclext:"filename,optional"`
	}

	ruleConfig.Filename = config.DefaultOutputsFileName

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	expected := fmt.Sprintf("%s.tf", ruleConfig.Filename)
	return enforceBlockFileBoundary(runner, r, expected, "output", 0)
}
