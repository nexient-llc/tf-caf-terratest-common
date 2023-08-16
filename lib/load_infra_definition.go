package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// aborts test if an error
func LoadRequestedInfraDefinitionFromTfVarsOut(t *testing.T, dir string, testFile string, testConfigOut any) {
	jsonOutPath := convertTFVars2JsonFile(t, dir, testFile)
	convertJsonFile2InfraDefinitionOut(t, jsonOutPath, testConfigOut)
	defer func() {
		os.Remove(jsonOutPath)
	}()

}

func convertJsonFile2InfraDefinitionOut(t *testing.T, jsonOutPath string, testConfigOut any) {

	bytes, err := os.ReadFile(jsonOutPath)
	require.NoError(t, err)
	err = json.Unmarshal(bytes, testConfigOut)
	require.NoError(t, err)

}

func convertTFVars2JsonFile(t *testing.T, dir string, testFile string) (jsonOutPath string) {
	const (
		tmp_json_filename = ".ignore_me.json"
	)
	jsonOutPath = filepath.Join(dir, tmp_json_filename)
	err := terraform.HCLFileToJSONFile(filepath.Join(dir, testFile), jsonOutPath)
	require.NoError(t, err)
	return jsonOutPath
}
