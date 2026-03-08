package rules

import (
	"fmt"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type VariablesFileRule struct {
	tflint.DefaultRule
}

func NewVariablesFileRule() *VariablesFileRule { return &VariablesFileRule{} }

func (r *VariablesFileRule) Name() string { return config.RulePrefix + "_variables_file" }

func (r *VariablesFileRule) Enabled() bool { return true }

func (r *VariablesFileRule) Severity() tflint.Severity { return severity() }

func (r *VariablesFileRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		Filename string `hclext:"filename,optional"`
	}

	ruleConfig.Filename = config.DefaultVariablesFileName

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	expected := fmt.Sprintf("%s.tf", ruleConfig.Filename)
	return enforceBlockFileBoundary(runner, r, expected, "variable")
}
