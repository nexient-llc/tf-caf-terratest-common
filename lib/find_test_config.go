package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FindTestConfig(testConfigFolderDefault string, infraTFVarFileNameDefault string) (string, string) {
	var infraTestConfigFolder, infraTestConfigVarFileName = testConfigFolderDefault, infraTFVarFileNameDefault
	if configEnvVar, envVarExists1 := os.LookupEnv("DSO_INFRA_TEST_CONFIG_FOLDER"); envVarExists1 {
		infraTestConfigFolder = configEnvVar
	}
	if testVarFileNameEnvVar, envVarExist2 := os.LookupEnv("DSO_INFRA_TEST_CONFIG_TFVAR_FILENAME"); envVarExist2 {
		infraTestConfigVarFileName = testVarFileNameEnvVar
	}
	return infraTestConfigFolder, infraTestConfigVarFileName
}
func IsExamplesFolder(t *testing.T, dir string) bool {
	return strings.HasSuffix(filepath.ToSlash(dir), "examples") || strings.HasSuffix(filepath.ToSlash(dir), "examples/")
}
func ListAllExamples(t *testing.T, examplesTFFolder string) []string {
	var folders []string
	examplesSubFolders, err := os.ReadDir(examplesTFFolder)
	assert.NoError(t, err)
	for _, anExampleFolder := range examplesSubFolders {
		if anExampleFolder.IsDir() {
			folders = append(folders, filepath.Join(examplesTFFolder, anExampleFolder.Name()))
		}
	}
	return folders
}
func IsSkipThisTestRequested(tfFoldr2Test string) bool {
	isTestDisabledCheck := map[string]bool{"false": true, "no": true, "n": true}
	envVarName := "DSO_INFRA_TEST_SKIP_TEST_" + filepath.Base(tfFoldr2Test)
	envVar, envVarExists := os.LookupEnv(envVarName)
	if envVarExists && !isTestDisabledCheck[strings.ToLower(envVar)] {
		fmt.Println("env var " + envVarName + " is set: skipping test for for " + tfFoldr2Test)
		return true
	} else {
		fmt.Println("env var " + envVarName + " is not set: executing test for for " + tfFoldr2Test)
		return false
	}
}
