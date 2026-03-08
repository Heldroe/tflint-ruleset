package rules

import (
    "github.com/Heldroe/tflint-ruleset/internal/config"
    "github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformBlockFileRule struct {
    tflint.DefaultRule

    SetupNumber   string `hcl:"setup_number,optional"`
    SetupFileName string `hcl:"setup_filename,optional"`
}

func NewTerraformBlockFileRule() *TerraformBlockFileRule {
    return &TerraformBlockFileRule{
        SetupNumber:   config.DefaultSetupNumber,
        SetupFileName: config.DefaultSetupFileName,
    }
}

func (r *TerraformBlockFileRule) Name() string {
    return config.RulePrefix + "_setup_file"
}

func (r *TerraformBlockFileRule) Enabled() bool {
    return true
}

func (r *TerraformBlockFileRule) Severity() tflint.Severity {
    return severity()
}

func (r *TerraformBlockFileRule) Check(runner tflint.Runner) error {

    expected := expectedFile(r.SetupNumber, r.SetupFileName)

    blocks, err := moduleBlocks(runner)
    if err != nil {
        return err
    }

    for _, block := range blocks {

        filename := block.DefRange.Filename

        if block.Type == "terraform" {

            if filename != expected {
                emit(runner, r,
                    "terraform block must be defined in "+expected,
                    block.DefRange,
                )
            }

            continue
        }

        if filename == expected {

            emit(runner, r,
                "only terraform blocks are allowed in "+expected,
                block.DefRange,
            )

        }
    }

    return nil
}
