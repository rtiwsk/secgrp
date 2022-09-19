package securitygroup

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ec2api interface {
	ModifyInstanceAttribute(ctx context.Context, params *ec2.ModifyInstanceAttributeInput, optFns ...func(*ec2.Options)) (*ec2.ModifyInstanceAttributeOutput, error)
	DescribeInstanceAttribute(ctx context.Context, params *ec2.DescribeInstanceAttributeInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceAttributeOutput, error)
}

type SG struct {
	client     ec2api
	instanceId string
}

func New(c aws.Config, instanceId string) *SG {
	return &SG{
		client:     ec2.NewFromConfig(c),
		instanceId: instanceId,
	}
}

func (s *SG) Add(ctx context.Context, groupId string) error {
	if groupId == "" {
		return errors.New("The security group is not specified.")
	}

	groups, err := s.List(ctx)
	if err != nil {
		return err
	}

	if existId(groups, groupId) {
		return errors.New("Already has the security group.")
	}

	groupIds := makeGroupIds(groups, groupId, "")

	if err = s.modify(ctx, groupIds); err != nil {
		return err
	}

	return nil
}

func (s *SG) Remove(ctx context.Context, groupId string) error {
	if groupId == "" {
		return errors.New("The security group is not specified.")
	}

	groups, err := s.List(ctx)
	if err != nil {
		return err
	}

	if !existId(groups, groupId) {
		return errors.New("The security group does not exist.")
	}

	groupIds := makeGroupIds(groups, "", groupId)

	if err = s.modify(ctx, groupIds); err != nil {
		return err
	}

	return nil
}

func (s *SG) List(ctx context.Context) ([]types.GroupIdentifier, error) {
	res, err := s.client.DescribeInstanceAttribute(ctx, &ec2.DescribeInstanceAttributeInput{
		Attribute:  types.InstanceAttributeNameGroupSet,
		InstanceId: aws.String(s.instanceId),
	})
	if err != nil {
		return nil, err
	}

	return res.Groups, nil
}

func (s *SG) modify(ctx context.Context, groupIds []string) error {
	_, err := s.client.ModifyInstanceAttribute(ctx, &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(s.instanceId),
		Groups:     groupIds,
	})
	if err != nil {
		return err
	}

	return nil
}

func existId(groups []types.GroupIdentifier, groupId string) bool {
	for _, g := range groups {
		if *g.GroupId == groupId {
			return true
		}
	}

	return false
}

func makeGroupIds(groups []types.GroupIdentifier, addGroupId, excludeGroupId string) []string {
	var groupIds []string
	for _, g := range groups {
		if *g.GroupId == excludeGroupId {
			continue
		}
		groupIds = append(groupIds, *g.GroupId)
	}

	if addGroupId != "" {
		groupIds = append(groupIds, addGroupId)
	}

	return groupIds
}
