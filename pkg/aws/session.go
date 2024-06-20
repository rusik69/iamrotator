package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rusik69/iamrotator/pkg/types"
)

// CreateSession creates a new AWS session
func CreateSession(cfg types.AWS) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
