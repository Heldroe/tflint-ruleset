package main

import "github.com/terraform-linters/tflint-plugin-sdk/plugin"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: NewRuleSet(),
	})
}
