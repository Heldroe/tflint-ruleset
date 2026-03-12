package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type NoBackendBlockRule struct {
	tflint.DefaultRule
}

func NewNoBackendBlockRule() *NoBackendBlockRule {
	return &NoBackendBlockRule{}
}

func (r *NoBackendBlockRule) Name() string {
	return config.RulePrefix + "_no_backend_block"
}

func (r *NoBackendBlockRule) Enabled() bool {
	return true
}

func (r *NoBackendBlockRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NoBackendBlockRule) Check(runner tflint.Runner) error {
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
			if block.Type == "terraform" {
				for _, innerBlock := range block.Body.Blocks {
					if innerBlock.Type == "backend" {
						runner.EmitIssue(
							r,
							"backend configuration is not allowed in the terraform block",
							innerBlock.TypeRange,
						)
					}
				}
			}
		}
	}

	return nil
}
