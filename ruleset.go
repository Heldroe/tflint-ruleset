package main

import (
    "github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
    "github.com/Heldroe/tflint-ruleset-terraform-style/rules"
    "github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func NewRuleSet() tflint.RuleSet {
    return &tflint.BuiltinRuleSet{
        Name:    config.RulePrefix,
        Version: "0.1.0",
        Rules: []tflint.Rule{
            rules.NewFilenamePatternRule(),

            rules.NewVariablesFileRule(),
            rules.NewOutputsFileRule(),
            rules.NewLocalsFileRule(),
            rules.NewDataFileRule(),
            rules.NewTerraformBlockFileRule(),

            rules.NewFileEndNewlineRule(),
            rules.NewBlockSpacingRule(),
        },
    }
}
