package rules

import (
	"bytes"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type MapAssignmentRule struct {
	tflint.DefaultRule
}

func NewMapAssignmentRule() *MapAssignmentRule {
	return &MapAssignmentRule{}
}

func (r *MapAssignmentRule) Name() string {
	return config.RulePrefix + "_map_assignment"
}

func (r *MapAssignmentRule) Enabled() bool {
	return true
}

func (r *MapAssignmentRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *MapAssignmentRule) Link() string {
	return ruleLink("map_assignment")
}

func (r *MapAssignmentRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		src := file.Bytes
		visitBodyExprs(body, func(expr hclsyntax.Expression) {
			if e, ok := expr.(*hclsyntax.ObjectConsExpr); ok {
				checkMapAssignment(runner, r, src, e)
			}
		})
	}

	return nil
}

func checkMapAssignment(runner tflint.Runner, rule tflint.Rule, src []byte, expr *hclsyntax.ObjectConsExpr) {
	for _, item := range expr.Items {
		keyEnd := item.KeyExpr.Range().End.Byte
		valStart := item.ValueExpr.Range().Start.Byte

		if keyEnd >= valStart || valStart > len(src) {
			continue
		}

		gap := src[keyEnd:valStart]
		if bytes.Contains(gap, []byte(":")) {
			runner.EmitIssue(
				rule,
				"maps must use the equal '=' sign to assign keys to values",
				item.KeyExpr.Range(),
			)
		}
	}
}
