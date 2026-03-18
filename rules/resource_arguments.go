package rules

import (
	"sort"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ResourceArgumentsRule struct {
	tflint.DefaultRule
}

func NewResourceArgumentsRule() *ResourceArgumentsRule {
	return &ResourceArgumentsRule{}
}

func (r *ResourceArgumentsRule) Name() string {
	return config.RulePrefix + "_resource_arguments"
}

func (r *ResourceArgumentsRule) Enabled() bool {
	return true
}

func (r *ResourceArgumentsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *ResourceArgumentsRule) Link() string {
	return ruleLink("resource_arguments")
}

func (r *ResourceArgumentsRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		walkBlocksResourceArguments(runner, r, file.Bytes, body)
	}

	return nil
}

func walkBlocksResourceArguments(runner tflint.Runner, rule tflint.Rule, src []byte, body *hclsyntax.Body) {
	for _, block := range body.Blocks {
		checkBlockArguments(runner, rule, src, block)
		walkBlocksResourceArguments(runner, rule, src, block.Body)
	}
}

type blockItem struct {
	Name      string
	Type      string // "attribute" or "block"
	Range     hcl.Range
	EmitRange hcl.Range
}

func checkBlockArguments(runner tflint.Runner, rule tflint.Rule, src []byte, block *hclsyntax.Block) {
	var items []blockItem
	for name, attr := range block.Body.Attributes {
		items = append(items, blockItem{
			Name:      name,
			Type:      "attribute",
			Range:     attr.Range(),
			EmitRange: attr.Range(),
		})
	}
	for _, b := range block.Body.Blocks {
		fullRange := hcl.Range{
			Filename: b.TypeRange.Filename,
			Start:    b.TypeRange.Start,
			End:      b.Body.SrcRange.End,
		}
		items = append(items, blockItem{
			Name:      b.Type,
			Type:      "block",
			Range:     fullRange,
			EmitRange: b.TypeRange,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Range.Start.Byte < items[j].Range.Start.Byte
	})

	if len(items) == 0 {
		return
	}

	firstIdx := 0
	for i, item := range items {
		if item.Name == "count" || item.Name == "for_each" {
			if i != firstIdx {
				isOk := false
				if i == 1 {
					prev := items[0]
					if prev.Name == "count" || prev.Name == "for_each" {
						isOk = true
						firstIdx = 1
					}
				}
				if !isOk {
					runner.EmitIssue(
						rule,
						"'count' and 'for_each' arguments must be the first thing declared in the block",
						item.EmitRange,
					)
				}
			} else {
				if i+1 < len(items) {
					checkEmptyLineAfter(runner, rule, src, item, items[i+1])
				}
				firstIdx++
			}
		}
	}

	if block.Type == "module" {
		for i, item := range items {
			if item.Name == "source" {
				for j := 0; j < i; j++ {
					prev := items[j]
					if prev.Name != "count" && prev.Name != "for_each" {
						runner.EmitIssue(
							rule,
							"'source' argument must be the first thing declared (but below 'count' and 'for_each')",
							item.EmitRange,
						)
						break
					}
				}
			}
		}
	}

	lastIdx := len(items) - 1
	hasDependsOn := false

	for i, item := range items {
		if item.Name == "depends_on" {
			hasDependsOn = true
			if i != lastIdx {
				runner.EmitIssue(
					rule,
					"'depends_on' must be the last thing declared in the block",
					item.EmitRange,
				)
			}
			if i > 0 {
				checkEmptyLineAbove(runner, rule, src, items[i-1], item)
			}
		}
	}

	for i, item := range items {
		if item.Name == "lifecycle" && item.Type == "block" {
			expectedIdx := lastIdx
			if hasDependsOn {
				expectedIdx = lastIdx - 1
			}
			if i != expectedIdx {
				runner.EmitIssue(
					rule,
					"'lifecycle' block must be the last thing declared in the block (but before 'depends_on')",
					item.EmitRange,
				)
			}
		}
	}
}

func checkEmptyLineAfter(runner tflint.Runner, rule tflint.Rule, src []byte, current, next blockItem) {
	start := current.Range.End.Byte
	end := next.Range.Start.Byte

	if start >= end || start >= len(src) {
		return
	}

	gap := src[start:end]
	newlines := strings.Count(string(gap), "\n")

	if newlines < 2 {
		runner.EmitIssue(
			rule,
			"there must be an empty blank line after "+current.Name,
			current.EmitRange,
		)
	}
}

func checkEmptyLineAbove(runner tflint.Runner, rule tflint.Rule, src []byte, prev, current blockItem) {
	start := prev.Range.End.Byte
	end := current.Range.Start.Byte

	if start >= end || start >= len(src) {
		return
	}

	gap := src[start:end]
	newlines := strings.Count(string(gap), "\n")

	if newlines < 2 {
		runner.EmitIssue(
			rule,
			"there must be an empty blank line above "+current.Name,
			current.EmitRange,
		)
	}
}
