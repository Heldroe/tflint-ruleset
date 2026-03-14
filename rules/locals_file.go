package rules

import (
    "fmt"

    "github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
    "github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type LocalsFileRule struct {
    tflint.DefaultRule
}

func NewLocalsFileRule() *LocalsFileRule { return &LocalsFileRule{} }

func (r *LocalsFileRule) Name() string { return config.RulePrefix + "_locals_file" }

func (r *LocalsFileRule) Enabled() bool { return true }

func (r *LocalsFileRule) Severity() tflint.Severity { return severity() }

func (r *LocalsFileRule) Link() string {
	return ruleLink("locals_file")
}

func (r *LocalsFileRule) Check(runner tflint.Runner) error {
    var ruleConfig struct {
        Filename string `hclext:"filename,optional"`
    }

    ruleConfig.Filename = config.DefaultLocalsFileName

    if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
        return err
    }

    expected := fmt.Sprintf("%s.tf", ruleConfig.Filename)
    return enforceBlockFileBoundary(runner, r, expected, "locals", 1)
    }

