package main

import (
	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/Heldroe/tflint-ruleset-terraform-style/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var version = "0.2.2"

func NewRuleSet() tflint.RuleSet {
	return &tflint.BuiltinRuleSet{
		Name:    config.RulePrefix,
		Version: version,
		Rules: []tflint.Rule{
			rules.NewFilenamePatternRule(),
			rules.NewVariablesFileRule(),
			rules.NewOutputsFileRule(),
			rules.NewLocalsFileRule(),
			rules.NewDataFileRule(),
			rules.NewTerraformBlockFileRule(),
			rules.NewResourceFileRule(),
			rules.NewNoBackendBlockRule(),
			rules.NewFileEndNewlineRule(),
			rules.NewBlockSpacingRule(),
			rules.NewNoEmptyFileRule(),
			rules.NewTrailingCommaRule(),
			rules.NewMapAssignmentRule(),
			rules.NewCommentStyleRule(),
			rules.NewResourceArgumentsRule(),
			rules.NewVariableArgumentsRule(),
			rules.NewOutputArgumentsRule(),
			rules.NewBlockInternalSpacingRule(),
			rules.NewNoProviderArgumentRule(),
			rules.NewStructureLayoutRule(),
		},
	}
}
