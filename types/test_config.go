package types

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type GenericTFModuleConfig struct {
	//the framework standard subset of attributes
	Naming_prefix      string            `json:"naming_prefix"`
	Environment        string            `json:"environment"`
	Environment_number string            `json:"environment_number"`
	Resource_number    string            `json:"resource_number"`
	Tags               map[string]string `json:"tags"`
	//to be extended by the TF module specific attrs
}

type TestContext struct {
	TestConfig                any // pointer to a TF module specific inheritance of GenericTFModuleConfig
	TestConfigFldrName        string
	TestConfigFileName        string
	TerratestTerraformOptions *terraform.Options
	CurrentTestName           string
}

func (ctx *TestContext) IsCurrentTest(testName string) bool {
	return ctx.CurrentTestName == testName
}
func (ctx *TestContext) EnabledOnlyForTests(t *testing.T, testName ...string) {
	for _, testName := range testName {
		if ctx.CurrentTestName == testName {
			return
		}
	}
	t.SkipNow()
}

type SecurityGroupT struct {
	EgressWithCidrBlocks []struct {
		CidrBlocksCommaSeparated string `json:"cidr_blocks"`
		CidrBlocks               []string
		FromPort                 int    `json:"from_port"`
		Protocol                 string `json:"protocol"`
		ToPort                   int    `json:"to_port"`
	} `json:"egress_with_cidr_blocks"`
	IngressCidrBlocks []string `json:"ingress_cidr_blocks"`
	IngressRules      []string `json:"ingress_rules"`
}
