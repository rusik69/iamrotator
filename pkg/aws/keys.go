package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/aws/session"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/rusik69/iamrotator/pkg/types"
)

// ListAccessKeys lists the access keys for all users
func ListAccessKeys(sess aws.Config, cfg types.AWSConfig) ([]types.AWSAccessKey, error) {
	var accessKeys []types.AWSAccessKey
	org := organizations.NewFromConfig(sess)
	err := org.ListAccounts(context.TODO(), &organizations.ListAccountsInput{},
		func(page *organizations.ListAccountsOutput, lastPage bool) bool {
			for _, account := range page.Accounts {
				creds := stscreds.NewCredentials(sess, fmt.Sprintf("arn:aws:iam::%s:role/IAMAccessRole", *account.Id))
				stsCfg := aws.Config{Credentials: creds}
				
	if err != nil {
		return nil, err
	}
	return accessKeys, nil
}
