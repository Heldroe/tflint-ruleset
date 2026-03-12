package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type NoProviderBlockRule struct {
	tflint.DefaultRule
}

func NewNoProviderBlockRule() *NoProviderBlockRule {
	return &NoProviderBlockRule{}
}

func (r *NoProviderBlockRule) Name() string {
	return config.RulePrefix + "_no_provider_block"
}

func (r *NoProviderBlockRule) Enabled() bool {
	return true
}

func (r *NoProviderBlockRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NoProviderBlockRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			if block.Type == "provider" {
				runner.EmitIssue(
					r,
					"provider blocks are not allowed",
					block.TypeRange,
				)
			}
		}
	}

	return nil
}
