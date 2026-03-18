package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type VariableArgumentsRule struct {
	tflint.DefaultRule
}

func NewVariableArgumentsRule() *VariableArgumentsRule {
	return &VariableArgumentsRule{}
}

func (r *VariableArgumentsRule) Name() string {
	return config.RulePrefix + "_variable_arguments"
}

func (r *VariableArgumentsRule) Enabled() bool {
	return true
}

func (r *VariableArgumentsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *VariableArgumentsRule) Link() string {
	return ruleLink("variable_arguments")
}

var defaultVariableOrder = []string{"type", "nullable", "sensitive", "ephemeral", "default", "description", "validation"}

func (r *VariableArgumentsRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		Order []string `hclext:"order,optional"`
	}
	ruleConfig.Order = defaultVariableOrder

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
			if block.Type != "variable" {
				continue
			}
			checkArgumentOrder(runner, r, block, orderIndex, ruleConfig.Order, "variable")
		}
	}

	return nil
}
