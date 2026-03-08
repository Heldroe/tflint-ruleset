package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func severity() tflint.Severity {
	return tflint.ERROR
}

func expectedFile(number string, name string) string {
	return fmt.Sprintf("%s-%s.tf", number, name)
}

func moduleBlocks(runner tflint.Runner) ([]*hclext.Block, error) {

	schema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "terraform"},
			{Type: "variable"},
			{Type: "locals"},
			{Type: "output"},
			{Type: "resource"},
			{Type: "data"},
			{Type: "module"},
			{Type: "provider"},
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
