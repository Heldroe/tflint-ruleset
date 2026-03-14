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
// 3. Multi-line structures must have their first element on a new line.
// 4. Commas must not be at the start of a line.
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
		var firstRange *hcl.Range
		if len(e.Exprs) > 0 {
			r := e.Exprs[0].Range()
			firstRange = &r
		}
		checkStructure(runner, rule, src, e.Range(), "]", firstRange)
		for _, item := range e.Exprs {
			walkExprStructure(runner, rule, src, item)
		}
	case *hclsyntax.ObjectConsExpr:
		var firstRange *hcl.Range
		if len(e.Items) > 0 {
			r := e.Items[0].KeyExpr.Range()
			firstRange = &r
		}
		checkStructure(runner, rule, src, e.Range(), "}", firstRange)
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
func checkStructure(runner tflint.Runner, rule tflint.Rule, src []byte, rng hcl.Range, bracketChar string, firstItemRange *hcl.Range) {
	startLine := rng.Start.Line
	endLine := rng.End.Line

	// 1. Single-line check
	if startLine == endLine {
		// Check for trailing comma
		endByte := rng.End.Byte
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
			return
		}

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

	// 2. Multi-line check: First element on its own line
	if firstItemRange != nil {
		if firstItemRange.Start.Line == startLine {
			runner.EmitIssue(
				rule,
				"for multi-line structures, the first element must be on a new line",
				hcl.Range{
					Filename: firstItemRange.Filename,
					Start:    firstItemRange.Start,
					End:      hcl.Pos{Line: firstItemRange.Start.Line, Column: firstItemRange.Start.Column + 1, Byte: firstItemRange.Start.Byte + 1},
				},
			)
		}
	}

	// 3. Multi-line check: Closing bracket on own line
	closeBracketIdx := -1
	for i := rng.End.Byte - 1; i >= rng.Start.Byte; i-- {
		if i < len(src) && src[i] == bracketChar[0] {
			closeBracketIdx = i
			break
		}
	}

	if closeBracketIdx != -1 {
		lineStart := 0
		lineNum := 1
		for i := 0; i < closeBracketIdx; i++ {
			if src[i] == '\n' {
				lineStart = i + 1
				lineNum++
			}
		}

		prefix := src[lineStart:closeBracketIdx]
		if strings.TrimSpace(string(prefix)) != "" {
			runner.EmitIssue(
				rule,
				"multi-line structures must have the closing bracket on its own line",
				hcl.Range{
					Filename: rng.Filename,
					Start:    hcl.Pos{Line: lineNum, Column: closeBracketIdx - lineStart + 1, Byte: closeBracketIdx},
					End:      hcl.Pos{Line: lineNum, Column: closeBracketIdx - lineStart + 2, Byte: closeBracketIdx + 1},
				},
			)
		}
	}

	// 4. Comma placement check
	if rng.End.Byte <= len(src) {
		text := src[rng.Start.Byte:rng.End.Byte]
		lines := bytes.Split(text, []byte("\n"))
		currentOffset := rng.Start.Byte
		currentLine := rng.Start.Line

		for i, line := range lines {
			if i > 0 {
				trimmed := bytes.TrimLeft(line, " \t")
				if len(trimmed) > 0 && trimmed[0] == ',' {
					col := len(line) - len(trimmed) + 1
					commaOffset := currentOffset + (len(line) - len(trimmed))
					runner.EmitIssue(
						rule,
						"commas must not be at the start of a line",
						hcl.Range{
							Filename: rng.Filename,
							Start:    hcl.Pos{Line: currentLine, Column: col, Byte: commaOffset},
							End:      hcl.Pos{Line: currentLine, Column: col + 1, Byte: commaOffset + 1},
						},
					)
					break 
				}
			}
			currentOffset += len(line) + 1 // +1 for newline
			currentLine++
		}
	}
}
