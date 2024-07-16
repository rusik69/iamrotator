package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rusik69/iamrotator/pkg/types"
	"github.com/sirupsen/logrus"
)

// ListAccessKeys lists the access keys for all users
func ListAccessKeys(sess aws.Config, cfg types.AWSConfig) ([]types.AWSAccessKey, error) {
	var accessKeys []types.AWSAccessKey
	org := organizations.NewFromConfig(sess)
	input := &organizations.ListAccountsForParentInput{
		ParentId:   &cfg.OUID,
		MaxResults: aws.Int32(20),
	}
	var nextToken *string

	for {
		if nextToken != nil {
			input.NextToken = nextToken
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
				logrus.Error(err)
				continue
			}
			sess, err := CreateSession(types.AWSConfig{
				Region:          cfg.Region,
				AccessKeyID:     *stsRes.Credentials.AccessKeyId,
				SecretAccessKey: *stsRes.Credentials.SecretAccessKey,
				SessionToken:    *stsRes.Credentials.SessionToken,
			})
			if err != nil {
				logrus.Error(err)
				continue
			}
			iamSvc := iam.NewFromConfig(sess)
			var userMarker *string
			for {
				usersOutput, err := iamSvc.ListUsers(context.TODO(), &iam.ListUsersInput{
					MaxItems: aws.Int32(20),
					Marker:   userMarker,
				})
				if err != nil {
					logrus.Error(err)
					break
				}
				for _, user := range usersOutput.Users {
					keysOutput, err := iamSvc.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
						UserName: user.UserName,
					})
					if err != nil {
						logrus.Error(err)
						continue
					}
					for _, key := range keysOutput.AccessKeyMetadata {
						accessKeys = append(accessKeys, types.AWSAccessKey{
							UserName:    *user.UserName,
							AccessKeyID: *key.AccessKeyId,
							AccountID:   *account.Id,
						})
					}
				}
				if usersOutput.Marker == nil {
					break
				}
				userMarker = usersOutput.Marker
			}
		}
		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}
	return accessKeys, nil
}
