package servicediscovery

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/stretchr/testify/require"
)

func GetAwsServiceDiscoveryClient(t *testing.T) *servicediscovery.Client {
	awsCfg, err1 := config.LoadDefaultConfig(context.TODO())
	require.NoError(t, err1, "retrieve AWS default config")
	client := servicediscovery.NewFromConfig(awsCfg)
	return client
}
