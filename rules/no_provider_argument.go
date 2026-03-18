package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

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

func (r *NoProviderArgumentRule) Link() string {
	return ruleLink("no_provider_argument")
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
	if attr, ok := body.Attributes["provider"]; ok {
		runner.EmitIssue(
			rule,
			"the \"provider\" argument is not allowed",
			attr.Range(),
		)
	}

	for _, block := range body.Blocks {
		walkBodyNoProviderArg(runner, rule, block.Body)
	}
}
