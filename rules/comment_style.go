package rules

import (
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

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

func (r *CommentStyleRule) Link() string {
	return ruleLink("comment_style")
}

func (r *CommentStyleRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for name, file := range files {
		tokens, diags := hclsyntax.LexConfig(file.Bytes, name, hcl.InitialPos)
		if diags.HasErrors() {
			continue
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

	if strings.HasPrefix(text, "//") || strings.HasPrefix(text, "/*") {
		runner.EmitIssue(
			rule,
			"comments are only allowed using '#' (number sign)",
			token.Range,
		)
		return
	}

	if !strings.HasPrefix(text, "#") {
		return
	}

	content := strings.TrimRight(text[1:], " \t\r\n")
	if len(content) == 0 {
		return
	}

	allHashes := true
	for _, ch := range content {
		if ch != '#' {
			allHashes = false
			break
		}
	}
	if allHashes {
		return
	}

	if !strings.HasPrefix(content, " ") {
		runner.EmitIssue(
			rule,
			"there must be a single space between the '#' and the beginning of the comment",
			token.Range,
		)
	}
}
