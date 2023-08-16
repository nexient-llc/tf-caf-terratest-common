package lambda

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/nexient-llc/tf-caf-terratest-common/lib/tags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLambdaInvokable(t *testing.T, awsApiLambdaClient *lambda.Client, functionName string) {
	//no direct http request - dont assume the test tool has access over network to the  infra being tested
	lambdaOut, err := awsApiLambdaClient.Invoke(context.TODO(), &lambda.InvokeInput{FunctionName: aws.String(functionName)})
	assert.NoError(t, err)
	assert.Equal(t, 200, int(lambdaOut.StatusCode))
}

func TestLambdaTags(t *testing.T, awsApiLambdaClient *lambda.Client, functionName string, expectedTags map[string]string) {
	lambdaOut := getAWSLambdaFunction(t, awsApiLambdaClient, functionName)
	expectedTags = tags.ExtendExpectedTagsByThoseAddedByFramework(expectedTags, functionName)
	tags.CompareExpectedTagsVsActual(t, expectedTags, lambdaOut.Tags)
}

func WaitForLambdaSpinUp(t *testing.T, awsApiLambdaClient *lambda.Client, functionName string) {
	const MAX_ATTEMPTS = 20
	const WAIT_SEC_PER_ATTEMPT = 6
	var counter int
	for {
		if counter >= MAX_ATTEMPTS {
			t.Fatalf("bad: time out while waiting for lambda %s to appear in this Cloud", functionName)
		}
		_, err := awsApiLambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{FunctionName: aws.String(functionName)})
		if err == nil {
			break
		}
		counter++
		logger.Log(t, "waiting until AWS ECS spins up container(s)")
		time.Sleep(WAIT_SEC_PER_ATTEMPT * time.Second)
	}
}

// aborts test run if err using "require".
func getAWSLambdaFunction(t *testing.T, awsApiLambdaClient *lambda.Client, functionName string) *lambda.GetFunctionOutput {
	//until github.com/gruntwork-io/terratest/modules/aws.getFunction adds support for AWS SDK v2
	out, err := awsApiLambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{FunctionName: aws.String(functionName)})
	require.NoError(t, err)
	return out
}

func GetAWSApiLambdaClient(t *testing.T) *lambda.Client {
	awsApiLambdaClient := lambda.NewFromConfig(GetAWSConfig(t))
	return awsApiLambdaClient
}
func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
