package rules

import (
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type BlockInternalSpacingRule struct {
	tflint.DefaultRule
}

func NewBlockInternalSpacingRule() *BlockInternalSpacingRule {
	return &BlockInternalSpacingRule{}
}

func (r *BlockInternalSpacingRule) Name() string {
	return config.RulePrefix + "_block_internal_spacing"
}

func (r *BlockInternalSpacingRule) Enabled() bool {
	return true
}

func (r *BlockInternalSpacingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *BlockInternalSpacingRule) Link() string {
	return ruleLink("block_internal_spacing")
}

func (r *BlockInternalSpacingRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		walkBlocksInternalSpacing(runner, r, file.Bytes, body)
	}

	return nil
}

func walkBlocksInternalSpacing(runner tflint.Runner, rule tflint.Rule, src []byte, body *hclsyntax.Body) {
	for _, block := range body.Blocks {
		checkBlockInternalSpacing(runner, rule, src, block)
		walkBlocksInternalSpacing(runner, rule, src, block.Body)
	}
}

func checkBlockInternalSpacing(runner tflint.Runner, rule tflint.Rule, src []byte, block *hclsyntax.Block) {
	start := block.Body.SrcRange.Start.Byte
	end := block.Body.SrcRange.End.Byte

	if start >= end || end > len(src) {
		return
	}

	contentStart := start
	if src[start] == '{' {
		contentStart++
	}
	contentEnd := end
	if src[end-1] == '}' {
		contentEnd--
	}

	if contentStart >= contentEnd {
		return
	}

	content := string(src[contentStart:contentEnd])
	lines := strings.Split(content, "\n")

	if len(lines) <= 1 {
		return
	}

	if strings.TrimSpace(lines[0]) == "" && len(lines) > 1 && strings.TrimSpace(lines[1]) == "" {
		runner.EmitIssue(
			rule,
			"there can't be a spacing empty line at the very top of the block",
			block.Body.SrcRange,
		)
	}

	lastIdx := len(lines) - 1
	if strings.TrimSpace(lines[lastIdx]) == "" && strings.TrimSpace(lines[lastIdx-1]) == "" {
		runner.EmitIssue(
			rule,
			"there can't be a spacing empty line at the very bottom of the block",
			block.Body.SrcRange,
		)
	}

	consecutiveBlank := 0
	for i, line := range lines {
		if i == 0 && strings.TrimSpace(line) == "" {
			continue
		}
		if i == lastIdx && strings.TrimSpace(line) == "" {
			continue
		}

		if strings.TrimSpace(line) == "" {
			consecutiveBlank++
			if consecutiveBlank > 1 {
				runner.EmitIssue(
					rule,
					"there cannot be more than 1 empty line of spacing in a row within blocks",
					block.Body.SrcRange,
				)
				break
			}
		} else {
			consecutiveBlank = 0
		}
	}
}
