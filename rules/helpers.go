package rules

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var filenamePattern = regexp.MustCompile(`^(\d{2})-[a-z0-9-]+\.tf$`)

var validTerraformBlocks = map[string]bool{
	"check":     true,
	"data":      true,
	"import":    true,
	"locals":    true,
	"module":    true,
	"moved":     true,
	"output":    true,
	"provider":  true,
	"removed":   true,
	"resource":  true,
	"terraform": true,
	"variable":  true,
}

func validateAllowedBlocks(blocks []string) error {
	for _, b := range blocks {
		if !validTerraformBlocks[b] {
			return fmt.Errorf("invalid block type %q in allowed_blocks; valid types are: check, data, import, locals, module, moved, output, provider, removed, resource, terraform, variable", b)
		}
	}
	return nil
}

var specialFileRules = []struct {
	Name            string
	DefaultFilename string
}{
	{Name: config.RulePrefix + "_variables_file", DefaultFilename: config.DefaultVariablesFileName},
	{Name: config.RulePrefix + "_terraform_file", DefaultFilename: config.DefaultTerraformFileName},
	{Name: config.RulePrefix + "_locals_file", DefaultFilename: config.DefaultLocalsFileName},
	{Name: config.RulePrefix + "_data_file", DefaultFilename: config.DefaultDataFileName},
	{Name: config.RulePrefix + "_outputs_file", DefaultFilename: config.DefaultOutputsFileName},
}

func resolveSpecialFiles(runner tflint.Runner) map[string]bool {
	files := make(map[string]bool, len(specialFileRules))
	for _, r := range specialFileRules {
		var cfg struct {
			Filename string `hclext:"filename,optional"`
		}
		cfg.Filename = r.DefaultFilename
		if err := runner.DecodeRuleConfig(r.Name, &cfg); err == nil && cfg.Filename != "" {
			files[cfg.Filename] = true
		} else {
			files[r.DefaultFilename] = true
		}
	}
	return files
}

func ruleLink(suffix string) string {
	return fmt.Sprintf("https://github.com/Heldroe/tflint-ruleset-terraform-style/blob/main/docs/rules/%s.md", suffix)
}

func enforceFileAllowedBlocks(runner tflint.Runner, rule tflint.Rule, targetFile string, allowedBlocks []string, maxBlocks map[string]int) error {
	if err := validateAllowedBlocks(allowedBlocks); err != nil {
		return err
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	allowed := make(map[string]bool, len(allowedBlocks))
	for _, b := range allowedBlocks {
		allowed[b] = true
	}

	blockCounts := make(map[string]int)
	for filename, file := range files {
		if filepath.Base(filename) != targetFile {
			continue
		}

		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, b := range body.Blocks {
			if !allowed[b.Type] {
				runner.EmitIssue(rule,
					fmt.Sprintf("only %s blocks are allowed in %s; found %s", strings.Join(allowedBlocks, ", "), targetFile, b.Type),
					b.TypeRange,
				)
				continue
			}
			blockCounts[b.Type]++
			if max, ok := maxBlocks[b.Type]; ok && blockCounts[b.Type] > max {
				runner.EmitIssue(rule,
					fmt.Sprintf("only %d %s block(s) allowed in %s; found multiple", max, b.Type, targetFile),
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
