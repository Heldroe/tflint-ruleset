package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ResourceFileRule struct {
	tflint.DefaultRule
}

func NewResourceFileRule() *ResourceFileRule {
	return &ResourceFileRule{}
}

func (r *ResourceFileRule) Name() string {
	return config.RulePrefix + "_resource_file"
}

func (r *ResourceFileRule) Enabled() bool {
	return true
}

func (r *ResourceFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *ResourceFileRule) Link() string {
	return ruleLink("resource_file")
}

func (r *ResourceFileRule) Check(runner tflint.Runner) error {
	var ruleConfig struct {
		AllowedBlocks []string `hclext:"allowed_blocks,optional"`
	}

	ruleConfig.AllowedBlocks = []string{"check", "module", "moved", "removed", "resource"}

	if err := runner.DecodeRuleConfig(r.Name(), &ruleConfig); err != nil {
		return err
	}

	if err := validateAllowedBlocks(ruleConfig.AllowedBlocks); err != nil {
		return err
	}

	specialFiles := resolveSpecialFiles(runner)

	allowed := make(map[string]bool, len(ruleConfig.AllowedBlocks))
	for _, b := range ruleConfig.AllowedBlocks {
		allowed[b] = true
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		baseName := filepath.Base(filename)
		nameWithoutExt := strings.TrimSuffix(baseName, ".tf")
		if specialFiles[nameWithoutExt] {
			continue
		}

		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, b := range body.Blocks {
			if !allowed[b.Type] {
				runner.EmitIssue(r,
					fmt.Sprintf("only %s blocks are allowed in resource files; found %s in %s", strings.Join(ruleConfig.AllowedBlocks, ", "), b.Type, baseName),
					b.TypeRange,
				)
			}
		}
	}

	return nil
}
