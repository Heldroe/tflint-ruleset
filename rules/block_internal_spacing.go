package rules

import (
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// BlockInternalSpacingRule checks spacing inside blocks.
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

		walkBodyInternalSpacing(runner, r, file.Bytes, body)
	}

	return nil
}

func walkBodyInternalSpacing(runner tflint.Runner, rule tflint.Rule, src []byte, body *hclsyntax.Body) {
	for _, block := range body.Blocks {
		checkBlockSpacing(runner, rule, src, block)
		walkBodyInternalSpacing(runner, rule, src, block.Body)
	}
}

func checkBlockSpacing(runner tflint.Runner, rule tflint.Rule, src []byte, block *hclsyntax.Block) {
	start := block.Body.SrcRange.Start.Byte
	end := block.Body.SrcRange.End.Byte

	if start >= end || end > len(src) {
		return
	}

	// We want the content inside.
	// If src[start] == '{', skip it.
	// If src[end-1] == '}', skip it.

	contentStart := start
	if src[start] == '{' {
		contentStart++
	}
	contentEnd := end
	if src[end-1] == '}' {
		contentEnd--
	}

	if contentStart >= contentEnd {
		return // Empty block or single line empty "{}"
	}

	content := string(src[contentStart:contentEnd])
	lines := strings.Split(content, "\n")
	
	// Rule: "no spacing empty line at the very top".
	// This applies to multi-line blocks.
	
	if len(lines) <= 1 {
		return // Single line block, spacing rules don't apply same way
	}
	
	// Check top
	// If lines[0] is empty (whitespace), look at lines[1].
	// If lines[1] is ALSO empty (whitespace), then we have a blank line at the top.
	if strings.TrimSpace(lines[0]) == "" {
		if len(lines) > 1 && strings.TrimSpace(lines[1]) == "" {
			runner.EmitIssue(
				rule,
				"there can't be a spacing empty line at the very top of the block",
				block.Body.SrcRange,
			)
		}
	}
	
	// Check bottom
	// If last line (lines[len-1]) is empty (whitespace) -> text before `}` on same line.
	// If lines[len-2] is ALSO empty -> blank line before `}`.
	if len(lines) > 1 {
		lastIdx := len(lines) - 1
		if strings.TrimSpace(lines[lastIdx]) == "" {
			if strings.TrimSpace(lines[lastIdx-1]) == "" {
				runner.EmitIssue(
					rule,
					"there can't be a spacing empty line at the very bottom of the block",
					block.Body.SrcRange,
				)
			}
		}
	}

	consecutiveBlank := 0
	for i, line := range lines {
		// Skip first line if it's the brace line (empty)
		if i == 0 && strings.TrimSpace(line) == "" {
			continue
		}
		// Skip last line if it's the brace line (empty)
		if i == len(lines)-1 && strings.TrimSpace(line) == "" {
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
				// Break after finding one violation to avoid noise
				break
			}
		} else {
			consecutiveBlank = 0
		}
	}
}
