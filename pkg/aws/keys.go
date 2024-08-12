package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
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
			newSess, err := CreateSessionWithRole(sess, cfg, *account.Id)
			if err != nil {
				logrus.Error(err)
				continue
			}
			iamSvc := iam.NewFromConfig(newSess)
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
						lastUsedOutput, err := iamSvc.GetAccessKeyLastUsed(context.TODO(), &iam.GetAccessKeyLastUsedInput{
							AccessKeyId: key.AccessKeyId,
						})
						if err != nil {
							logrus.Error(err)
							continue
						}
						lastUsedDate := time.Time{}
						if lastUsedOutput.AccessKeyLastUsed.LastUsedDate != nil {
							lastUsedDate = *lastUsedOutput.AccessKeyLastUsed.LastUsedDate
						}
						accessKeys = append(accessKeys, types.AWSAccessKey{
							UserName:     *user.UserName,
							AccessKeyID:  *key.AccessKeyId,
							AccountID:    *account.Id,
							AccountName:  *account.Name,
							CreateDate:   *key.CreateDate,
							LastUsedDate: lastUsedDate,
							Status:       string(key.Status),
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

// RemoveAccessKey removes the access key for the user
func RemoveAccessKey(sess aws.Config, cfg types.AWSConfig, key types.AWSAccessKey) error {
	newSess, err := CreateSessionWithRole(sess, cfg, key.AccountID)
	if err != nil {
		return err
	}
	iamSvc := iam.NewFromConfig(newSess)
	input2 := &iam.DeleteAccessKeyInput{
		AccessKeyId: &key.AccessKeyID,
		UserName:    &key.UserName,
	}
	_, err = iamSvc.DeleteAccessKey(context.TODO(), input2)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccessKey creates a new access key for the user
func CreateAccessKey(sess aws.Config, cfg types.AWSConfig, user string) (types.AWSAccessKey, error) {
	newSess, err := CreateSessionWithRole(sess, cfg, cfg.AccountID)
	if err != nil {
		return types.AWSAccessKey{}, err
	}
	iamSvc := iam.NewFromConfig(newSess)
	input := &iam.CreateAccessKeyInput{
		UserName: &user,
	}
	result, err := iamSvc.CreateAccessKey(context.TODO(), input)
	if err != nil {
		return types.AWSAccessKey{}, err
	}
	return types.AWSAccessKey{
		UserName:        user,
		AccessKeyID:     *result.AccessKey.AccessKeyId,
		SecretAccessKey: *result.AccessKey.SecretAccessKey,
		AccountID:       cfg.AccountID,
	}, nil
}
