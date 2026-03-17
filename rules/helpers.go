package rules

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var filenamePattern = regexp.MustCompile(`^(\d{2})-[a-z0-9-]+\.tf$`)

func ruleLink(suffix string) string {
	return fmt.Sprintf("https://github.com/Heldroe/tflint-ruleset-terraform-style/blob/main/docs/rules/%s.md", suffix)
}

func enforceBlockFileBoundary(runner tflint.Runner, rule tflint.Rule, expectedFile, blockType string, maxBlocks int) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	blockCount := 0
	for filename, file := range files {
		baseName := filepath.Base(filename)

		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, b := range body.Blocks {
			if b.Type == blockType {
				if baseName != expectedFile {
					runner.EmitIssue(rule,
						fmt.Sprintf("%s blocks must be defined in %s", blockType, expectedFile),
						b.TypeRange,
					)
				} else {
					blockCount++
					if maxBlocks > 0 && blockCount > maxBlocks {
						runner.EmitIssue(rule,
							fmt.Sprintf("only %d %s block(s) allowed in %s; found multiple", maxBlocks, blockType, expectedFile),
							b.TypeRange,
						)
					}
				}
			}

			if baseName == expectedFile && b.Type != blockType {
				runner.EmitIssue(rule,
					fmt.Sprintf("only %s blocks are allowed in %s; found %s", blockType, expectedFile, b.Type),
					b.TypeRange,
				)
			}
		}
	}

	return nil
}

func visitBodyExprs(body *hclsyntax.Body, visit func(hclsyntax.Expression)) {
	for _, attr := range body.Attributes {
		visitExprTree(attr.Expr, visit)
	}
	for _, block := range body.Blocks {
		visitBodyExprs(block.Body, visit)
	}
}

func visitExprTree(expr hclsyntax.Expression, visit func(hclsyntax.Expression)) {
	visit(expr)
	switch e := expr.(type) {
	case *hclsyntax.TupleConsExpr:
		for _, item := range e.Exprs {
			visitExprTree(item, visit)
		}
	case *hclsyntax.ObjectConsExpr:
		for _, item := range e.Items {
			visitExprTree(item.KeyExpr, visit)
			visitExprTree(item.ValueExpr, visit)
		}
	case *hclsyntax.FunctionCallExpr:
		for _, arg := range e.Args {
			visitExprTree(arg, visit)
		}
	case *hclsyntax.ConditionalExpr:
		visitExprTree(e.Condition, visit)
		visitExprTree(e.TrueResult, visit)
		visitExprTree(e.FalseResult, visit)
	case *hclsyntax.ForExpr:
		visitExprTree(e.CollExpr, visit)
		if e.KeyExpr != nil {
			visitExprTree(e.KeyExpr, visit)
		}
		visitExprTree(e.ValExpr, visit)
		if e.CondExpr != nil {
			visitExprTree(e.CondExpr, visit)
		}
	case *hclsyntax.ParenthesesExpr:
		visitExprTree(e.Expression, visit)
	case *hclsyntax.BinaryOpExpr:
		visitExprTree(e.LHS, visit)
		visitExprTree(e.RHS, visit)
	case *hclsyntax.UnaryOpExpr:
		visitExprTree(e.Val, visit)
	case *hclsyntax.IndexExpr:
		visitExprTree(e.Collection, visit)
		visitExprTree(e.Key, visit)
	case *hclsyntax.SplatExpr:
		visitExprTree(e.Source, visit)
	case *hclsyntax.TemplateExpr:
		for _, part := range e.Parts {
			visitExprTree(part, visit)
		}
	case *hclsyntax.TemplateWrapExpr:
		visitExprTree(e.Wrapped, visit)
	}
}
