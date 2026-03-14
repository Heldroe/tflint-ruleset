package rules

import (
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CommentStyleRule enforces comment style: only '#' allowed, single space required.
type CommentStyleRule struct {
	tflint.DefaultRule
}

func NewCommentStyleRule() *CommentStyleRule {
	return &CommentStyleRule{}
}

func (r *CommentStyleRule) Name() string {
	return config.RulePrefix + "_comment_style"
}

func (r *CommentStyleRule) Enabled() bool {
	return true
}

func (r *CommentStyleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *CommentStyleRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for name, file := range files {
		tokens, diags := hclsyntax.LexConfig(file.Bytes, name, hcl.InitialPos)
		if diags.HasErrors() {
			continue // Skip files that don't parse cleanly (TFLint might handle them elsewhere)
		}

		for _, token := range tokens {
			if token.Type == hclsyntax.TokenComment {
				checkComment(runner, r, token)
			}
		}
	}

	return nil
}

func checkComment(runner tflint.Runner, rule tflint.Rule, token hclsyntax.Token) {
	text := string(token.Bytes)
	
	// Rule 1: Only '#' allowed
	if strings.HasPrefix(text, "//") || strings.HasPrefix(text, "/*") {
		runner.EmitIssue(
			rule,
			"comments are only allowed using '#' (number sign)",
			token.Range,
		)
		return
	}

	if strings.HasPrefix(text, "#") {
		// Rule 2: Single space required after '#', unless all '#'
		content := text[1:] // Strip leading '#'
		
		// Trim trailing whitespace (newline etc)
		content = strings.TrimRight(content, " \t\r\n")
		
		if len(content) == 0 {
			return
		}
		
		// Check if it's all '#'
		allHashes := true
		for _, r := range content {
			if r != '#' {
				allHashes = false
				break
			}
		}
		if allHashes {
			return
		}

		// Check for space
		if !strings.HasPrefix(content, " ") {
			runner.EmitIssue(
				rule,
				"there must be a single space between the '#' and the beginning of the comment",
				token.Range,
			)
		}
	}
}
