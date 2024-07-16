package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rusik69/iamrotator/pkg/types"
)

// CreateSession creates a new AWS session
func CreateSession(cfg types.AWSConfig) (aws.Config, error) {
	// Load the default configuration
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)),
	)
	if err != nil {
		return aws.Config{}, err
	}
	return awsCfg, nil
}

// CreateSessionWithRole creates a new AWS session with a role
func CreateSessionWithRole(sess aws.Config, cfg types.AWSConfig, accountID string) (aws.Config, error) {
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::" + accountID + ":role/" + cfg.RoleName),
		RoleSessionName: aws.String(cfg.RoleName),
	}
	stsSvc := sts.NewFromConfig(sess)
	stsRes, err := stsSvc.AssumeRole(context.TODO(), input)
	if err != nil {
		return aws.Config{}, err
	}
	newSess, err := CreateSession(types.AWSConfig{
		Region:          cfg.Region,
		AccessKeyID:     *stsRes.Credentials.AccessKeyId,
		SecretAccessKey: *stsRes.Credentials.SecretAccessKey,
		SessionToken:    *stsRes.Credentials.SessionToken,
	})
	return newSess, err
}
