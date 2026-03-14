package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// NoProviderArgumentRule forbids "provider" argument in blocks.
type NoProviderArgumentRule struct {
	tflint.DefaultRule
}

func NewNoProviderArgumentRule() *NoProviderArgumentRule {
	return &NoProviderArgumentRule{}
}

func (r *NoProviderArgumentRule) Name() string {
	return config.RulePrefix + "_no_provider_argument"
}

func (r *NoProviderArgumentRule) Enabled() bool {
	return false
}

func (r *NoProviderArgumentRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NoProviderArgumentRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		walkBodyNoProviderArg(runner, r, body)
	}

	return nil
}

func walkBodyNoProviderArg(runner tflint.Runner, rule tflint.Rule, body *hclsyntax.Body) {
	// Check attributes in current body
	if attr, ok := body.Attributes["provider"]; ok {
		runner.EmitIssue(
			rule,
			"the \"provider\" argument in not allowed",
			attr.Range(),
		)
	}

	// Recurse into blocks
	for _, block := range body.Blocks {
		walkBodyNoProviderArg(runner, rule, block.Body)
	}
}
