package rules

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func severity() tflint.Severity {
	return tflint.ERROR
}

func enforceBlockFileBoundary(runner tflint.Runner, rule tflint.Rule, expectedFile string, blockType string, maxBlocks int) error {
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
			// Rule 1: blockType blocks MUST be in the expected file
			if b.Type == blockType {
				if baseName != expectedFile {
					emit(runner, rule,
						fmt.Sprintf("%s blocks must be defined in %s", blockType, expectedFile),
						b.TypeRange,
					)
				} else {
					blockCount++
					if maxBlocks > 0 && blockCount > maxBlocks {
						emit(runner, rule,
							fmt.Sprintf("only %d %s block(s) allowed in %s; found multiple", maxBlocks, blockType, expectedFile),
							b.TypeRange,
						)
					}
				}
			}

			// Rule 2: The expected file MUST ONLY contain blockType blocks
			if baseName == expectedFile && b.Type != blockType {
				emit(runner, rule,
					fmt.Sprintf("only %s blocks are allowed in %s; found %s", blockType, expectedFile, b.Type),
					b.TypeRange,
				)
			}
		}
	}

	return nil
}

func moduleBlocks(runner tflint.Runner) ([]*hclext.Block, error) {
	schema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "terraform"},
			{Type: "locals"},
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "output", LabelNames: []string{"name"}},
			{Type: "module", LabelNames: []string{"name"}},
			{Type: "provider", LabelNames: []string{"name"}},
			{Type: "resource", LabelNames: []string{"type", "name"}},
			{Type: "data", LabelNames: []string{"type", "name"}},
		},
	}

	content, err := runner.GetModuleContent(schema, nil)
	if err != nil {
		return nil, err
	}

	return content.Blocks, nil
}

func emit(runner tflint.Runner, rule tflint.Rule, msg string, r hcl.Range) {
	runner.EmitIssue(rule, msg, r)
}
