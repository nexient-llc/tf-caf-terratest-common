package ec2

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nexient-llc/tf-caf-terratest-common/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func CompareSecurityGroupExpectedVsActual(expected types.SecurityGroupT, actual ec2_types.SecurityGroup) bool {
	isFound := false
	for _, expectedEgressRule := range expected.EgressWithCidrBlocks {
		for _, actualEgressRule := range actual.IpPermissionsEgress {
			if compareSecurityGroupPort(expectedEgressRule.FromPort, actualEgressRule.FromPort) &&
				compareSecurityGroupPort(expectedEgressRule.ToPort, actualEgressRule.ToPort) &&
				compareSecurityGroupIPProtocol(expectedEgressRule.Protocol, actualEgressRule.IpProtocol) &&
				compareSecurityGroupCidrsArrays(expectedEgressRule.CidrBlocks, actualEgressRule.IpRanges) {
				isFound = true
				break
			}
		}
	}
	return isFound
}

func compareSecurityGroupPort(expected int, actual *int32) bool {
	if (expected == 0) && (actual == nil) {
		return true
	}
	if expected == int(*actual) {
		return true
	}
	return false
}
func compareSecurityGroupIPProtocol(expected string, actual *string) bool {
	if (expected == "-1") && (actual == nil) {
		return true
	}
	if strings.EqualFold(expected, *actual) {
		return true
	}
	return false
}

func compareSecurityGroupCidrsArrays(requestedCidrBlocks []string, actualCidrBlocks []ec2_types.IpRange) bool {
	sort.Strings(requestedCidrBlocks)
	var ipRanges []string
	for _, iprange := range actualCidrBlocks {
		ipRanges = append(ipRanges, *iprange.CidrIp)
	}
	sort.Strings(ipRanges)
	//"0.0.0.0/0" has the same meaning as "" in AWS sec groups API
	if strings.EqualFold(requestedCidrBlocks[0], "0.0.0.0/0") && len(actualCidrBlocks) == 0 {
		return true
	}
	return slices.Equal(ipRanges, requestedCidrBlocks)
}

func GetAWSSecGroupByID(t *testing.T, secGrpId string, ec2ApiClient *ec2.Client) *ec2.DescribeSecurityGroupsOutput {
	input := ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{secGrpId},
	}
	secGroups, err := ec2ApiClient.DescribeSecurityGroups(context.Background(), &input)
	assert.NoError(t, err)
	assert.True(t, len(secGroups.SecurityGroups) == 1, "one sec group should exists for ID "+secGrpId)
	return secGroups
}

func GetAWSApiClientEC2(t *testing.T) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)
	return ec2.NewFromConfig(cfg)
}
