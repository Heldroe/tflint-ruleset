package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TrailingCommaRule checks for trailing commas in lists and maps.
type TrailingCommaRule struct {
	tflint.DefaultRule
}

func NewTrailingCommaRule() *TrailingCommaRule {
	return &TrailingCommaRule{}
}

func (r *TrailingCommaRule) Name() string {
	return config.RulePrefix + "_trailing_comma"
}

func (r *TrailingCommaRule) Enabled() bool {
	return true
}

func (r *TrailingCommaRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TrailingCommaRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		walkBody(runner, r, file.Bytes, body)
	}

	return nil
}

func walkBody(runner tflint.Runner, rule tflint.Rule, src []byte, body *hclsyntax.Body) {
	for _, attr := range body.Attributes {
		walkExpr(runner, rule, src, attr.Expr)
	}
	for _, block := range body.Blocks {
		walkBody(runner, rule, src, block.Body)
	}
}

func walkExpr(runner tflint.Runner, rule tflint.Rule, src []byte, expr hclsyntax.Expression) {
	switch e := expr.(type) {
	case *hclsyntax.TupleConsExpr: // List: [ ... ]
		checkList(runner, rule, src, e)
		for _, item := range e.Exprs {
			walkExpr(runner, rule, src, item)
		}
	case *hclsyntax.ObjectConsExpr: // Map/Object: { ... }
		checkMap(runner, rule, src, e)
		for _, item := range e.Items {
			walkExpr(runner, rule, src, item.KeyExpr)
			walkExpr(runner, rule, src, item.ValueExpr)
		}
	case *hclsyntax.FunctionCallExpr:
		for _, arg := range e.Args {
			walkExpr(runner, rule, src, arg)
		}
	case *hclsyntax.ConditionalExpr:
		walkExpr(runner, rule, src, e.Condition)
		walkExpr(runner, rule, src, e.TrueResult)
		walkExpr(runner, rule, src, e.FalseResult)
	case *hclsyntax.ForExpr:
		walkExpr(runner, rule, src, e.CollExpr)
		if e.KeyExpr != nil {
			walkExpr(runner, rule, src, e.KeyExpr)
		}
		walkExpr(runner, rule, src, e.ValExpr)
		if e.CondExpr != nil {
			walkExpr(runner, rule, src, e.CondExpr)
		}
	case *hclsyntax.ParenthesesExpr:
		walkExpr(runner, rule, src, e.Expression)
	case *hclsyntax.BinaryOpExpr:
		walkExpr(runner, rule, src, e.LHS)
		walkExpr(runner, rule, src, e.RHS)
	case *hclsyntax.UnaryOpExpr:
		walkExpr(runner, rule, src, e.Val)
	case *hclsyntax.IndexExpr:
		walkExpr(runner, rule, src, e.Collection)
		walkExpr(runner, rule, src, e.Key)
	case *hclsyntax.SplatExpr:
		walkExpr(runner, rule, src, e.Source)
		// e.Each is implicit usually, check?
	case *hclsyntax.TemplateExpr:
		for _, part := range e.Parts {
			walkExpr(runner, rule, src, part)
		}
	case *hclsyntax.TemplateWrapExpr:
		walkExpr(runner, rule, src, e.Wrapped)
	case *hclsyntax.AnonSymbolExpr:
		// Terminal
	case *hclsyntax.LiteralValueExpr:
		// Terminal
	case *hclsyntax.ScopeTraversalExpr:
		// Terminal
	case *hclsyntax.RelativeTraversalExpr:
		// Terminal
	default:
		// Traverse other expressions if they contain nested structures
	}
}

func checkList(runner tflint.Runner, rule tflint.Rule, src []byte, expr *hclsyntax.TupleConsExpr) {
	if len(expr.Exprs) == 0 {
		return
	}
	
	rng := expr.Range()
	startLine := rng.Start.Line
	endLine := rng.End.Line 

	if startLine == endLine {
		return // Single line, ignore
	}

	// Multi-line list
	// Check if the last item has a trailing comma
	lastItem := expr.Exprs[len(expr.Exprs)-1]
	lastItemEnd := lastItem.Range().End.Byte
	
	exprEnd := rng.End.Byte

	if lastItemEnd >= exprEnd || lastItemEnd >= len(src) || exprEnd > len(src) {
		return
	}

	// Extract text between last item and closing bracket
	gap := src[lastItemEnd:exprEnd]
	
	// Check for comma, ignoring comments
	hasComma := false
	for i := 0; i < len(gap); i++ {
		b := gap[i]
		if b == ',' {
			hasComma = true
			break
		}
		if b == '#' || (b == '/' && i+1 < len(gap) && (gap[i+1] == '/' || gap[i+1] == '*')) {
			// Comment start, stop looking
			break
		}
	}

	if !hasComma {
		runner.EmitIssue(
			rule,
			"lists defined in multiple lines must have a trailing comma on the last line",
			lastItem.Range(),
		)
	}
}

func checkMap(runner tflint.Runner, rule tflint.Rule, src []byte, expr *hclsyntax.ObjectConsExpr) {
	if len(expr.Items) == 0 {
		return
	}

	rng := expr.Range()
	startLine := rng.Start.Line
	endLine := rng.End.Line

	if startLine == endLine {
		return // Single line
	}

	// Multi-line map: NO trailing commas on ANY line.
	// We need to check after EACH item.
	for _, item := range expr.Items {
		// Check gap after ValueExpr
		valEnd := item.ValueExpr.Range().End.Byte
		
		limit := len(src)
		scannerIdx := int(valEnd)
		foundComma := false
		
		for scannerIdx < limit {
			b := src[scannerIdx]
			if b == '\n' || b == '\r' {
				break
			}
			if b == ',' {
				foundComma = true
				break
			}
			if b == '#' || (b == '/' && scannerIdx+1 < limit && (src[scannerIdx+1] == '/' || src[scannerIdx+1] == '*')) {
				break
			}
			scannerIdx++
		}

		if foundComma {
			runner.EmitIssue(
				rule,
				"maps defined in multiple lines must not have any trailing comma on any lines",
				hcl.Range{
					Filename: item.ValueExpr.Range().Filename,
					Start:    item.ValueExpr.Range().End,
					End:      hcl.Pos{Line: item.ValueExpr.Range().End.Line, Column: item.ValueExpr.Range().End.Column + 1, Byte: int(valEnd) + 1},
				},
			)
		}
	}
}
