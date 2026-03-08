package rules

import (
    "regexp"

    "github.com/hashicorp/hcl/v2"

    "github.com/Heldroe/tflint-ruleset-terraform-style/internal/config"
    "github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type FilenamePatternRule struct {
    tflint.DefaultRule
}

func NewFilenamePatternRule() *FilenamePatternRule {
    return &FilenamePatternRule{}
}

func (r *FilenamePatternRule) Name() string {
    return config.RulePrefix + "_filename_pattern"
}

func (r *FilenamePatternRule) Enabled() bool {
    return true
}

func (r *FilenamePatternRule) Severity() tflint.Severity {
    return tflint.ERROR
}

func (r *FilenamePatternRule) Check(runner tflint.Runner) error {

    files, err := runner.GetFiles()
    if err != nil {
        return err
    }

    re := regexp.MustCompile(`^\d{2}-[a-z0-9-]+\.tf$`)

    for name := range files {

        if !re.MatchString(name) {

            runner.EmitIssue(
                r,
                "terraform file name must match XX-name.tf",
                hcl.Range{
                    Filename: name,
                    Start: hcl.Pos{
                        Line:   1,
                        Column: 1,
                    },
                    End: hcl.Pos{
                        Line:   1,
                        Column: 1,
                    },
                },
            )

        }
    }

    return nil
}
