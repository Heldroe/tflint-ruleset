package rules

import (
	"bytes"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

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

func (r *StructureLayoutRule) Link() string {
	return ruleLink("structure_layout")
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

		src := file.Bytes
		visitBodyExprs(body, func(expr hclsyntax.Expression) {
			switch e := expr.(type) {
			case *hclsyntax.TupleConsExpr:
				var firstRange *hcl.Range
				if len(e.Exprs) > 0 {
					rng := e.Exprs[0].Range()
					firstRange = &rng
				}
				checkStructure(runner, r, src, e.Range(), "]", firstRange)
			case *hclsyntax.ObjectConsExpr:
				var firstRange *hcl.Range
				if len(e.Items) > 0 {
					rng := e.Items[0].KeyExpr.Range()
					firstRange = &rng
				}
				checkStructure(runner, r, src, e.Range(), "}", firstRange)
			}
		})
	}

	return nil
}

func checkStructure(runner tflint.Runner, rule tflint.Rule, src []byte, rng hcl.Range, bracketChar string, firstItemRange *hcl.Range) {
	startLine := rng.Start.Line
	endLine := rng.End.Line

	if startLine == endLine {
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

	if firstItemRange != nil && firstItemRange.Start.Line == startLine {
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
			currentOffset += len(line) + 1
			currentLine++
		}
	}
}
