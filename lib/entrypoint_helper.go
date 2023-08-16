package lib

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func NewTerratestTerraformOptions(dir string, testFile string) *terraform.Options {
	terraformOptions := &terraform.Options{
		TerraformDir: dir,
		VarFiles:     []string{testFile},
		NoColor:      true,
		Logger:       logger.Discard,
	}
	return terraformOptions
}
