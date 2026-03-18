package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type OutputArgumentsRule struct {
	tflint.DefaultRule
}

func NewOutputArgumentsRule() *OutputArgumentsRule {
	return &OutputArgumentsRule{}
}

func (r *OutputArgumentsRule) Name() string {
	return config.RulePrefix + "_output_arguments"
}

func (r *OutputArgumentsRule) Enabled() bool {
	return true
}

func (r *OutputArgumentsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *OutputArgumentsRule) Link() string {
	return ruleLink("output_arguments")
}

var defaultOutputOrder = []string{"description", "sensitive", "ephemeral", "value", "precondition", "depends_on"}

func (r *OutputArgumentsRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		Order []string `hclext:"order,optional"`
	}
	ruleConfig.Order = defaultOutputOrder

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	orderIndex := make(map[string]int, len(ruleConfig.Order))
	for i, name := range ruleConfig.Order {
		orderIndex[name] = i
	}

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
			if block.Type != "output" {
				continue
			}
			checkArgumentOrder(runner, r, block, orderIndex, ruleConfig.Order, "output")
		}
	}

	return nil
}
