package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/nexient-llc/tf-caf-terratest-common/types"
)

const (
	FUNC_NAME_CONVENTION_COMPOSABLE_PREFIX = "TestComposable"
)

type TestFunc func(t *testing.T, testCtx types.TestContext)

func internalRunSetupTestTeardown(t *testing.T, dir string, testFile string, testCtx types.TestContext, targetInfraReadOnly bool, testFunc ...TestFunc) {
	// check if test should be skipped
	if IsSkipThisTestRequested(dir) {
		return
	}
	// load config
	LoadRequestedInfraDefinitionFromTfVarsOut(t, dir, testFile, testCtx.TestConfig)
	testCtx.TerratestTerraformOptions = NewTerratestTerraformOptions(dir, testFile)
	testCtx.CurrentTestName = filepath.Base(dir)

	// Apply - unless Regression Mode or local iteration Stage skipped
	if !targetInfraReadOnly {
		stage := test_structure.RunTestStage

		defer stage(t, "teardown_test_"+testCtx.CurrentTestName, func() {
			teardownTestGeneric(t, testCtx)
		})
		stage(t, "setup_test_"+testCtx.CurrentTestName, func() {
			setupTestGeneric(t, testCtx)
		})
	}

	// get IDs of deployed cloud resources
	// Verify health of cloud resources. When healthy - test passed
	for _, fun := range testFunc {
		fun(t, testCtx)
	}
	// Teardown - unless Regression Mode or local iteration Stage skipped

}

func setupTestGeneric(t *testing.T, testCtx types.TestContext) {
	terraform.InitAndApplyAndIdempotent(t, testCtx.TerratestTerraformOptions)
}

func teardownTestGeneric(t *testing.T, testCtx types.TestContext) {
	terraform.Destroy(t, testCtx.TerratestTerraformOptions)
}

func ForEveryExampleRunTest(t *testing.T, infraTests2Run []string, testVarFileName string, testCtx types.TestContext, testFunc ...TestFunc) {
	for _, tfFoldr2Test := range infraTests2Run {
		internalRunSetupTestTeardownDestructive(t, tfFoldr2Test, testVarFileName, testCtx, testFunc...)
	}
}

func internalRunSetupTestTeardownReadonly(t *testing.T, dir string, testFile string, testCtx types.TestContext, testFunc ...TestFunc) {
	internalRunSetupTestTeardown(t, dir, testFile, testCtx, true, testFunc...)
}

func internalRunSetupTestTeardownDestructive(t *testing.T, dir string, testFile string, testCtx types.TestContext, testFunc ...TestFunc) {
	internalRunSetupTestTeardown(t, dir, testFile, testCtx, false, testFunc...)
}

func RunSetupTestTeardown(t *testing.T, dir string, testFile string, testCtx types.TestContext, testFunc ...TestFunc) {
	infra2TestTfFolder, testVarFileName := FindTestConfig(dir, testFile)
	var infraTests2Run []string
	if IsExamplesFolder(t, infra2TestTfFolder) {
		infraTests2Run = append(infraTests2Run, ListAllExamples(t, infra2TestTfFolder)...)
	} else {
		infraTests2Run = append(infraTests2Run, infra2TestTfFolder)
	}
	ForEveryExampleRunTest(t, infraTests2Run, testVarFileName, testCtx, testFunc...)

}

func RunNonDestructiveTest(t *testing.T, dir string, testFile string, testCtx types.TestContext, testFunc ...TestFunc) {
	infra2TestTfFolder, testVarFileName := FindTestConfig(dir, testFile)
	demandDeployedTerraformStateExists(t, infra2TestTfFolder)
	demandAllTests2RunAreComposableOnes(t, testFunc)
	internalRunSetupTestTeardownReadonly(t, infra2TestTfFolder, testVarFileName, testCtx, testFunc...)
}

// stops test by 'require' if error
func demandDeployedTerraformStateExists(t *testing.T, infra2TestTfFolder string) {
	if !isDeployedTerraformStateDetected(infra2TestTfFolder) {
		fmt.Println("bad:can not detect deployed TF folder. dir .terraform does not exists or unreadable")
		t.FailNow()
	}
}
func isDeployedTerraformStateDetected(infra2TestTfFolder string) bool {
	exists := false
	dotTerraformDir := filepath.Join(infra2TestTfFolder, ".terraform")
	fileInfo, err := os.Stat(dotTerraformDir)
	if err == nil {
		exists = fileInfo.IsDir()
	}
	return exists
}

// stops test by 'require' if wrong user input
func demandAllTests2RunAreComposableOnes(t *testing.T, testFunc []TestFunc) {
	for _, fun := range testFunc {
		funcName := getFunctionNameShort(fun)
		if !isFunctionNameAComposableOne(funcName) {
			fmt.Println("bad: for high env regression, only low level/composable test code can be used and no text fixtures.\n" +
				"test function " + funcName + " does not look like composable one - does not follow naming convention for those.\n" +
				"should have name prefix " + FUNC_NAME_CONVENTION_COMPOSABLE_PREFIX)
			t.FailNow()
		}
	}
}

func getFunctionNameShort(funPtr interface{}) string {
	nameParts := strings.Split(runtime.FuncForPC(reflect.ValueOf(funPtr).Pointer()).Name(), ".")
	return nameParts[len(nameParts)-1]
}

func isFunctionNameAComposableOne(name string) bool {
	return strings.HasPrefix(name, FUNC_NAME_CONVENTION_COMPOSABLE_PREFIX)
}
