package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type NoEmptyFileRule struct {
	tflint.DefaultRule
}

func NewNoEmptyFileRule() *NoEmptyFileRule {
	return &NoEmptyFileRule{}
}

func (r *NoEmptyFileRule) Name() string {
	return config.RulePrefix + "_no_empty_file"
}

func (r *NoEmptyFileRule) Enabled() bool {
	return true
}

func (r *NoEmptyFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NoEmptyFileRule) Link() string {
	return ruleLink("no_empty_file")
}

func (r *NoEmptyFileRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		if len(body.Blocks) == 0 && len(body.Attributes) == 0 {
			runner.EmitIssue(r, "file cannot be empty", body.MissingItemRange())
		}
	}

	return nil
}
