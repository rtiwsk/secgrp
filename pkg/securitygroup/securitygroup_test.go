package securitygroup

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

var testgroups = []types.GroupIdentifier{
	{
		GroupId:   aws.String("sg-1234567890abcdefg"),
		GroupName: aws.String("testGroup1"),
	},
	{
		GroupId:   aws.String("sg-1234567890hijklmn"),
		GroupName: aws.String("testGroup2"),
	},
}

func Test_existId(t *testing.T) {
	tests := []struct {
		name    string
		groups  []types.GroupIdentifier
		groupId string
		want    bool
	}{
		{
			name:    "exist the groupId",
			groups:  testgroups,
			groupId: "sg-1234567890abcdefg",
			want:    true,
		},
		{
			name:    "does not exist the groupId",
			groups:  testgroups,
			groupId: "sg-1234567890opqrstu",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := existId(tt.groups, tt.groupId)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_makeGroupIds(t *testing.T) {
	tests := []struct {
		name           string
		groups         []types.GroupIdentifier
		addGroupId     string
		excludeGroupId string
		want           []string
	}{
		{
			name:       "add the groupId",
			groups:     testgroups,
			addGroupId: "sg-1234567890opqrstu",
			want: []string{
				"sg-1234567890abcdefg",
				"sg-1234567890hijklmn",
				"sg-1234567890opqrstu",
			},
		},
		{

			name:           "exclude the groupId",
			groups:         testgroups,
			excludeGroupId: "sg-1234567890hijklmn",
			want: []string{
				"sg-1234567890abcdefg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeGroupIds(tt.groups, tt.addGroupId, tt.excludeGroupId)
			assert.Equal(t, got, tt.want)
		})
	}
}
