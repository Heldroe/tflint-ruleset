package rules

import (
	"strings"

	"github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type FileEndNewlineRule struct{ tflint.DefaultRule }

func NewFileEndNewlineRule() *FileEndNewlineRule { return &FileEndNewlineRule{} }

func (r *FileEndNewlineRule) Name() string {
	return config.RulePrefix + "_file_end_newline"
}

func (r *FileEndNewlineRule) Enabled() bool { return true }

func (r *FileEndNewlineRule) Severity() tflint.Severity { return severity() }

func (r *FileEndNewlineRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for name, file := range files {
		text := string(file.Bytes)

		// Fail if it doesn't end in a newline, OR if it ends in two or more newlines
		if !strings.HasSuffix(text, "\n") || strings.HasSuffix(text, "\n\n") {
			emit(runner, r,
				"file must end with exactly one newline",
				hcl.Range{
					Filename: name,
					Start:    hcl.Pos{Line: 1, Column: 1},
					End:      hcl.Pos{Line: 1, Column: 1},
				},
			)
		}
	}

	return nil
}
