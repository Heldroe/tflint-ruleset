package rules

import (
	"sort"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type BlockSpacingRule struct {
	tflint.DefaultRule
}

func NewBlockSpacingRule() *BlockSpacingRule {
	return &BlockSpacingRule{}
}

func (r *BlockSpacingRule) Name() string {
	return config.RulePrefix + "_block_spacing"
}

func (r *BlockSpacingRule) Enabled() bool {
	return true
}

func (r *BlockSpacingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *BlockSpacingRule) Link() string {
	return ruleLink("block_spacing")
}

func (r *BlockSpacingRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		blocks := body.Blocks
		if len(blocks) < 2 {
			continue
		}

		sort.Slice(blocks, func(i, j int) bool {
			return blocks[i].TypeRange.Start.Byte < blocks[j].TypeRange.Start.Byte
		})

		for i := 1; i < len(blocks); i++ {
			prev := blocks[i-1]
			curr := blocks[i]

			startByte := prev.Body.SrcRange.End.Byte
			endByte := curr.TypeRange.Start.Byte

			if startByte >= endByte || startByte >= len(file.Bytes) || endByte > len(file.Bytes) {
				continue
			}

			gap := string(file.Bytes[startByte:endByte])
			lines := strings.Split(gap, "\n")

			blankLineCount := 0
			maxConsecutiveBlanks := 0
			consecutiveBlanks := 0
			hasNonCommentContent := false
			if len(lines) > 2 {
				for _, line := range lines[1 : len(lines)-1] {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" {
						blankLineCount++
						consecutiveBlanks++
						if consecutiveBlanks > maxConsecutiveBlanks {
							maxConsecutiveBlanks = consecutiveBlanks
						}
					} else {
						consecutiveBlanks = 0
						if !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") {
							hasNonCommentContent = true
						}
					}
				}
			}

			if blankLineCount < 1 || maxConsecutiveBlanks > 1 || hasNonCommentContent {
				runner.EmitIssue(r,
					"blocks must be separated by exactly one blank line",
					curr.TypeRange,
				)
			}
		}
	}

	return nil
}
