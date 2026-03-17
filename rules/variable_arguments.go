package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
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
			checkVariableOrder(runner, r, block, orderIndex, ruleConfig.Order)
		}
	}

	return nil
}

type orderedItem struct {
	name      string
	emitRange hcl.Range
	sortKey   int
}

func checkVariableOrder(runner tflint.Runner, rule tflint.Rule, block *hclsyntax.Block, orderIndex map[string]int, order []string) {
	var items []orderedItem
	for name, attr := range block.Body.Attributes {
		if idx, ok := orderIndex[name]; ok {
			items = append(items, orderedItem{name: name, emitRange: attr.Range(), sortKey: idx})
		}
	}
	for _, b := range block.Body.Blocks {
		if idx, ok := orderIndex[b.Type]; ok {
			items = append(items, orderedItem{name: b.Type, emitRange: b.TypeRange, sortKey: idx})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].emitRange.Start.Byte < items[j].emitRange.Start.Byte
	})

	maxSeen := -1
	for _, it := range items {
		if it.sortKey < maxSeen {
			var expectedBefore string
			for _, prev := range items {
				if prev.sortKey == maxSeen {
					expectedBefore = prev.name
					break
				}
			}
			runner.EmitIssue(rule,
				fmt.Sprintf("'%s' must be declared before '%s' in variable blocks (expected order: %s)",
					it.name, expectedBefore, formatPresentOrder(order, items)),
				it.emitRange,
			)
		}
		if it.sortKey > maxSeen {
			maxSeen = it.sortKey
		}
	}
}

func formatPresentOrder(order []string, items []orderedItem) string {
	present := make(map[string]bool, len(items))
	for _, it := range items {
		present[it.name] = true
	}
	var parts []string
	for _, name := range order {
		if present[name] {
			parts = append(parts, name)
		}
	}
	return strings.Join(parts, ", ")
}
