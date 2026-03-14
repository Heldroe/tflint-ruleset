package rules

import (
	"bytes"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StructureLayoutRule enforces layout rules for lists and maps.
// 1. Single-line structures must not have a trailing comma.
// 2. Multi-line structures must have the closing bracket on its own line.
// 3. Commas must not be at the start of a line.
type StructureLayoutRule struct {
	tflint.DefaultRule
}

func NewStructureLayoutRule() *StructureLayoutRule {
	return &StructureLayoutRule{}
}

func (r *StructureLayoutRule) Name() string {
	return config.RulePrefix + "_structure_layout"
}

func (r *StructureLayoutRule) Enabled() bool {
	return true
}

func (r *StructureLayoutRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *StructureLayoutRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		walkBodyStructure(runner, r, file.Bytes, body)
	}

	return nil
}

func walkBodyStructure(runner tflint.Runner, rule tflint.Rule, src []byte, body *hclsyntax.Body) {
	for _, attr := range body.Attributes {
		walkExprStructure(runner, rule, src, attr.Expr)
	}
	for _, block := range body.Blocks {
		walkBodyStructure(runner, rule, src, block.Body)
	}
}

func walkExprStructure(runner tflint.Runner, rule tflint.Rule, src []byte, expr hclsyntax.Expression) {
	switch e := expr.(type) {
	case *hclsyntax.TupleConsExpr:
		checkStructure(runner, rule, src, e.Range(), "]", e.Exprs)
		for _, item := range e.Exprs {
			walkExprStructure(runner, rule, src, item)
		}
	case *hclsyntax.ObjectConsExpr:
		checkStructure(runner, rule, src, e.Range(), "}", nil)
		for _, item := range e.Items {
			walkExprStructure(runner, rule, src, item.KeyExpr)
			walkExprStructure(runner, rule, src, item.ValueExpr)
		}
	// Recursive traversal for other expressions
	case *hclsyntax.FunctionCallExpr:
		for _, arg := range e.Args {
			walkExprStructure(runner, rule, src, arg)
		}
	case *hclsyntax.ConditionalExpr:
		walkExprStructure(runner, rule, src, e.Condition)
		walkExprStructure(runner, rule, src, e.TrueResult)
		walkExprStructure(runner, rule, src, e.FalseResult)
	case *hclsyntax.ForExpr:
		walkExprStructure(runner, rule, src, e.CollExpr)
		if e.KeyExpr != nil {
			walkExprStructure(runner, rule, src, e.KeyExpr)
		}
		walkExprStructure(runner, rule, src, e.ValExpr)
		if e.CondExpr != nil {
			walkExprStructure(runner, rule, src, e.CondExpr)
		}
	case *hclsyntax.ParenthesesExpr:
		walkExprStructure(runner, rule, src, e.Expression)
	case *hclsyntax.BinaryOpExpr:
		walkExprStructure(runner, rule, src, e.LHS)
		walkExprStructure(runner, rule, src, e.RHS)
	case *hclsyntax.UnaryOpExpr:
		walkExprStructure(runner, rule, src, e.Val)
	case *hclsyntax.IndexExpr:
		walkExprStructure(runner, rule, src, e.Collection)
		walkExprStructure(runner, rule, src, e.Key)
	case *hclsyntax.SplatExpr:
		walkExprStructure(runner, rule, src, e.Source)
	case *hclsyntax.TemplateExpr:
		for _, part := range e.Parts {
			walkExprStructure(runner, rule, src, part)
		}
	case *hclsyntax.TemplateWrapExpr:
		walkExprStructure(runner, rule, src, e.Wrapped)
	}
}

// checkStructure checks layout for a generic list/map expression
// bracketChar is "]" or "}"
func checkStructure(runner tflint.Runner, rule tflint.Rule, src []byte, rng hcl.Range, bracketChar string, listItems []hclsyntax.Expression) {
	startLine := rng.Start.Line
	endLine := rng.End.Line

	// 1. Single-line check
	if startLine == endLine {
		// Check for trailing comma
		// Find the closing bracket position
		// Range.End is usually after the closing bracket.
		// Scan backwards from End.Byte - 1
		endByte := rng.End.Byte
		
		// Safety check
		if endByte > len(src) {
			return
		}

		scanIdx := endByte - 1
		for scanIdx >= rng.Start.Byte {
			if scanIdx < len(src) && src[scanIdx] == bracketChar[0] {
				break
			}
			scanIdx--
		}

		if scanIdx < rng.Start.Byte {
			return // Should not happen
		}

		// Now scan backwards from before bracket for comma
		foundComma := false
		for i := scanIdx - 1; i >= rng.Start.Byte; i-- {
			b := src[i]
			if b == ' ' || b == '\t' {
				continue
			}
			if b == ',' {
				foundComma = true
				break
			}
			// Hit content (not whitespace/comma)
			break
		}

		if foundComma {
			runner.EmitIssue(
				rule,
				"single-line structures must not have a trailing comma",
				rng,
			)
		}
		return
	}

	closeBracketIdx := -1
	for i := rng.End.Byte - 1; i >= rng.Start.Byte; i-- {
		if i < len(src) && src[i] == bracketChar[0] {
			closeBracketIdx = i
			break
		}
	}

	if closeBracketIdx != -1 {
		lineStart := 0
		for i := closeBracketIdx - 1; i >= 0; i-- {
			if src[i] == '\n' {
				lineStart = i + 1
				break
			}
		}

		// Check text from lineStart to closeBracketIdx
		prefix := src[lineStart:closeBracketIdx]
		if strings.TrimSpace(string(prefix)) != "" {
			runner.EmitIssue(
				rule,
				"multi-line structures must have the closing bracket on its own line",
				rng,
			)
		}
	}

	// Comma placement check
	// Scan the raw text of the expression for comma at start of line
	// Range: rng.Start.Byte to rng.End.Byte
	if rng.End.Byte <= len(src) {
		text := src[rng.Start.Byte:rng.End.Byte]
		
		lines := bytes.Split(text, []byte("\n"))
		// Skip first line (it contains `[` or `start`).
		for i := 1; i < len(lines); i++ {
			line := lines[i]
			trimmed := bytes.TrimSpace(line)
			if len(trimmed) > 0 && trimmed[0] == ',' {
				runner.EmitIssue(
					rule,
					"commas must not be at the start of a line",
					rng,
				)
				break 
			}
		}
	}
}
