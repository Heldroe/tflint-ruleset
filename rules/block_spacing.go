package rules

import (
    "sort"
    "strings"

    "github.com/hashicorp/hcl/v2/hclsyntax"
    "github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
    "github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type BlockSpacingRule struct{ tflint.DefaultRule }

func NewBlockSpacingRule() *BlockSpacingRule { return &BlockSpacingRule{} }

func (r *BlockSpacingRule) Name() string {
    return config.RulePrefix + "_block_spacing"
}

func (r *BlockSpacingRule) Enabled() bool { return true }

func (r *BlockSpacingRule) Severity() tflint.Severity { return severity() }

func (r *BlockSpacingRule) Link() string {
	return ruleLink("block_spacing")
}

func (r *BlockSpacingRule) Check(runner tflint.Runner) error {
    files, err := runner.GetFiles()
    if err != nil {
        return err
    }

    for _, file := range files {
        // hclext abstracts away physical braces. By casting the raw file body
        // to hclsyntax.Body, we get access to exact brace byte boundaries.
        body, ok := file.Body.(*hclsyntax.Body)
        if !ok {
            continue
        }

        // Grab all top-level blocks in this file
        blocks := body.Blocks
        if len(blocks) < 2 {
            continue
        }

        // Sort blocks by their starting byte just to be safe
        sort.Slice(blocks, func(i, j int) bool {
            return blocks[i].TypeRange.Start.Byte < blocks[j].TypeRange.Start.Byte
        })

        for i := 1; i < len(blocks); i++ {
            prev := blocks[i-1]
            curr := blocks[i]

            // prev.Body.SrcRange covers everything inside { ... }
            // SrcRange.End is the exact position right after the closing '}'
            startByte := prev.Body.SrcRange.End.Byte

            // curr.TypeRange.Start is the exact position of the next block's keyword
            endByte := curr.TypeRange.Start.Byte

            // Safety check
            if startByte >= endByte || startByte >= len(file.Bytes) || endByte > len(file.Bytes) {
                continue
            }

            // Extract the literal text gap between the closing '}' and next block
            gap := string(file.Bytes[startByte:endByte])
            lines := strings.Split(gap, "\n")

            blankLineCount := 0
            if len(lines) > 2 {
                // We skip lines[0] (trailing text on the previous closing brace's line)
                // and lines[len(lines)-1] (indentation before the current block)
                for _, line := range lines[1 : len(lines)-1] {
                    if strings.TrimSpace(line) == "" {
                        blankLineCount++
                    }
                }
            }

            if blankLineCount != 1 {
                // We can just pass curr.TypeRange to your emit function.
                // It acts exactly like DefRange but is pulled directly from the raw AST.
                emit(runner, r,
                    "blocks must be separated by exactly one blank line",
                    curr.TypeRange,
                )
            }
        }
    }

    return nil
}
