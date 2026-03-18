package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

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

func (r *TrailingCommaRule) Link() string {
	return ruleLink("trailing_comma")
}

func (r *TrailingCommaRule) Check(runner tflint.Runner) error {
	ruleConfig := struct {
		ExcludeSingleElement bool `hclext:"exclude_single_element,optional"`
	}{
		ExcludeSingleElement: false,
	}

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
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

		src := file.Bytes
		exclude := ruleConfig.ExcludeSingleElement
		visitBodyExprs(body, func(expr hclsyntax.Expression) {
			switch e := expr.(type) {
			case *hclsyntax.TupleConsExpr:
				checkList(runner, r, src, e, exclude)
			case *hclsyntax.ObjectConsExpr:
				checkMap(runner, r, src, e, exclude)
			}
		})
	}

	return nil
}

func checkList(runner tflint.Runner, rule tflint.Rule, src []byte, expr *hclsyntax.TupleConsExpr, excludeSingle bool) {
	if len(expr.Exprs) == 0 {
		return
	}

	if excludeSingle && len(expr.Exprs) == 1 {
		return
	}

	rng := expr.Range()
	if rng.Start.Line == rng.End.Line {
		return
	}

	lastItem := expr.Exprs[len(expr.Exprs)-1]
	lastItemEnd := lastItem.Range().End.Byte
	exprEnd := rng.End.Byte

	if lastItemEnd >= exprEnd || lastItemEnd >= len(src) || exprEnd > len(src) {
		return
	}

	gap := src[lastItemEnd:exprEnd]

	hasComma := false
	for i := 0; i < len(gap); i++ {
		b := gap[i]
		if b == ',' {
			hasComma = true
			break
		}
		if b == '#' || (b == '/' && i+1 < len(gap) && (gap[i+1] == '/' || gap[i+1] == '*')) {
			break
		}
	}

	if !hasComma {
		lastItemRng := lastItem.Range()
		startPos := lastItemRng.End
		if startPos.Column > 1 {
			startPos.Column--
			startPos.Byte--
		}
		runner.EmitIssue(
			rule,
			"lists defined in multiple lines must have a trailing comma on the last line",
			hcl.Range{
				Filename: lastItemRng.Filename,
				Start:    startPos,
				End:      lastItemRng.End,
			},
		)
	}
}

func checkMap(runner tflint.Runner, rule tflint.Rule, src []byte, expr *hclsyntax.ObjectConsExpr, excludeSingle bool) {
	if len(expr.Items) == 0 {
		return
	}

	if excludeSingle && len(expr.Items) == 1 {
		return
	}

	rng := expr.Range()
	if rng.Start.Line == rng.End.Line {
		return
	}

	for _, item := range expr.Items {
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
