package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rusik69/iamrotator/pkg/types"
)

// ListAccessKeys lists the access keys for all users
func ListAccessKeys(sess aws.Config, cfg types.AWSConfig) ([]types.AWSAccessKey, error) {
	var accessKeys []types.AWSAccessKey
	org := organizations.NewFromConfig(sess)
	input := &organizations.ListAccountsForParentInput{
		ParentId: &cfg.OUID,
	}
	result, err := org.ListAccountsForParent(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	for _, account := range result.Accounts {
		stsSvc := sts.NewFromConfig(sess)
		input := &sts.AssumeRoleInput{
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", *account.Id, cfg.RoleName)),
			RoleSessionName: aws.String(cfg.RoleName),
		}
		stsRes, err := stsSvc.AssumeRole(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		iamSvc := iam.NewFromConfig(aws.Config{
			Credentials: iam.CredentialsProvider{
				Source: stsRes.Credentials,
			},
		})
		usersOutput, err := iamSvc.ListUsers(context.TODO(), &iam.ListUsersInput{})
		if err != nil {
			return nil, err
		}
		for _, user := range usersOutput.Users {
			keysOutput, err := iamSvc.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
				UserName: user.UserName,
			})
			if err != nil {
				return nil, err
			}
			for _, key := range keysOutput.AccessKeyMetadata {
				accessKeys = append(accessKeys, types.AWSAccessKey{
					AccountID:   *account.Id,
					UserName:    *user.UserName,
					AccessKeyID: *key.AccessKeyId,
				})
			}
		}
	}
	return accessKeys, nil
}
