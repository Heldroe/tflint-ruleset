package rules

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// NoEmptyFileRule checks if a file is empty or contains no blocks/attributes.
type NoEmptyFileRule struct {
	tflint.DefaultRule
}

// NewNoEmptyFileRule returns a new NoEmptyFileRule instance.
func NewNoEmptyFileRule() *NoEmptyFileRule {
	return &NoEmptyFileRule{}
}

func (r *NoEmptyFileRule) Name() string {
	return config.RulePrefix + "_no_empty_file"
}

func (r *NoEmptyFileRule) Enabled() bool {
	return true
}

func (r *NoEmptyFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NoEmptyFileRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		if len(body.Blocks) == 0 && len(body.Attributes) == 0 {
			// Get the file range if possible, or create a dummy range
			rng := body.MissingItemRange()
			runner.EmitIssue(
				r,
				"file cannot be empty",
				rng,
			)
		}
	}

	return nil
}
