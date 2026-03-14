package rules

import (
	"bytes"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MapAssignmentRule checks that maps use '=' instead of ':'.
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

		walkBodyMapAssignment(runner, r, file.Bytes, body)
	}

	return nil
}

func walkBodyMapAssignment(runner tflint.Runner, rule tflint.Rule, src []byte, body *hclsyntax.Body) {
	for _, attr := range body.Attributes {
		walkExprMapAssignment(runner, rule, src, attr.Expr)
	}
	for _, block := range body.Blocks {
		walkBodyMapAssignment(runner, rule, src, block.Body)
	}
}

func walkExprMapAssignment(runner tflint.Runner, rule tflint.Rule, src []byte, expr hclsyntax.Expression) {
	switch e := expr.(type) {
	case *hclsyntax.TupleConsExpr:
		for _, item := range e.Exprs {
			walkExprMapAssignment(runner, rule, src, item)
		}
	case *hclsyntax.ObjectConsExpr:
		checkMapAssignment(runner, rule, src, e)
		for _, item := range e.Items {
			walkExprMapAssignment(runner, rule, src, item.KeyExpr)
			walkExprMapAssignment(runner, rule, src, item.ValueExpr)
		}
	case *hclsyntax.FunctionCallExpr:
		for _, arg := range e.Args {
			walkExprMapAssignment(runner, rule, src, arg)
		}
	case *hclsyntax.ConditionalExpr:
		walkExprMapAssignment(runner, rule, src, e.Condition)
		walkExprMapAssignment(runner, rule, src, e.TrueResult)
		walkExprMapAssignment(runner, rule, src, e.FalseResult)
	case *hclsyntax.ForExpr:
		walkExprMapAssignment(runner, rule, src, e.CollExpr)
		if e.KeyExpr != nil {
			walkExprMapAssignment(runner, rule, src, e.KeyExpr)
		}
		walkExprMapAssignment(runner, rule, src, e.ValExpr)
		if e.CondExpr != nil {
			walkExprMapAssignment(runner, rule, src, e.CondExpr)
		}
	case *hclsyntax.ParenthesesExpr:
		walkExprMapAssignment(runner, rule, src, e.Expression)
	case *hclsyntax.BinaryOpExpr:
		walkExprMapAssignment(runner, rule, src, e.LHS)
		walkExprMapAssignment(runner, rule, src, e.RHS)
	case *hclsyntax.UnaryOpExpr:
		walkExprMapAssignment(runner, rule, src, e.Val)
	case *hclsyntax.IndexExpr:
		walkExprMapAssignment(runner, rule, src, e.Collection)
		walkExprMapAssignment(runner, rule, src, e.Key)
	case *hclsyntax.SplatExpr:
		walkExprMapAssignment(runner, rule, src, e.Source)
	case *hclsyntax.TemplateExpr:
		for _, part := range e.Parts {
			walkExprMapAssignment(runner, rule, src, part)
		}
	case *hclsyntax.TemplateWrapExpr:
		walkExprMapAssignment(runner, rule, src, e.Wrapped)
	}
}

func checkMapAssignment(runner tflint.Runner, rule tflint.Rule, src []byte, expr *hclsyntax.ObjectConsExpr) {
	for _, item := range expr.Items {
		keyEnd := item.KeyExpr.Range().End.Byte
		valStart := item.ValueExpr.Range().Start.Byte

		if keyEnd >= valStart || valStart > len(src) {
			continue
		}

		// Text between key and value
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
