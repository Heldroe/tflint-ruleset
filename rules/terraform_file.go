package rules

import (
	"fmt"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformBlockFileRule struct {
	tflint.DefaultRule
}

func NewTerraformBlockFileRule() *TerraformBlockFileRule {
	return &TerraformBlockFileRule{}
}

func (r *TerraformBlockFileRule) Name() string {
	return config.RulePrefix + "_terraform_file"
}

func (r *TerraformBlockFileRule) Enabled() bool {
	return true
}

func (r *TerraformBlockFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformBlockFileRule) Link() string {
	return ruleLink("terraform_file")
}

func (r *TerraformBlockFileRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		Filename      string   `hclext:"filename,optional"`
		AllowedBlocks []string `hclext:"allowed_blocks,optional"`
	}

	ruleConfig.Filename = config.DefaultTerraformFileName
	ruleConfig.AllowedBlocks = []string{"terraform"}

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	expected := fmt.Sprintf("%s.tf", ruleConfig.Filename)
	return enforceFileAllowedBlocks(runner, r, expected, ruleConfig.AllowedBlocks, map[string]int{"terraform": 1})
}
