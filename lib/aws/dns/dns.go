package dns

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53_types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/require"
)

func GetAWSApiRoute53Client(t *testing.T) *route53.Client {
	awsApiLambdaClient := route53.NewFromConfig(GetAWSConfig(t))
	return awsApiLambdaClient
}
func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
func NameNormalize(name string) string {
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	return name
}

func GetHostedZoneById(t *testing.T, zone_id string) *route53.GetHostedZoneOutput {
	zone, err := GetAWSApiRoute53Client(t).GetHostedZone(context.Background(), &route53.GetHostedZoneInput{
		Id: &zone_id,
	})
	require.NoError(t, err)
	return zone
}
func LookupDNSRecordInPublicRoute53ZoneByDNSProtocol(t *testing.T, zone_id string, fullQualifiedRecordName string, recType string) (*route53.TestDNSAnswerOutput, error) {
	return GetAWSApiRoute53Client(t).TestDNSAnswer(context.Background(), &route53.TestDNSAnswerInput{
		HostedZoneId: &zone_id,
		RecordName:   &fullQualifiedRecordName,
		RecordType:   route53_types.RRType(recType),
	})
}
