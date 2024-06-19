package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// ListIamUsers lists all IAM users
func ListIamUsers(sess *session.Session) ([]string, error) {
	svc := iam.New(sess)
	result, err := svc.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}
	var userNames []string
	for _, user := range result.Users {
		userNames = append(userNames, *user.UserName)
	}
	return userNames, nil
}

// CreateIamUser creates a new IAM user
func CreateIamUser(sess *session.Session, userName string) error {
	svc := iam.New(sess)
	input := &iam.CreateUserInput{
		UserName: &userName,
	}
	_, err := svc.CreateUser(input)
	if err != nil {
		return err
	}
	policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": "sts:AssumeRole",
                "Resource": "*"
            }
        ]
    }`
	putUserPolicyInput := &iam.PutUserPolicyInput{
		UserName:       aws.String(userName),
		PolicyName:     aws.String("AssumeRolePolicy"),
		PolicyDocument: aws.String(policy),
	}
	_, err = svc.PutUserPolicy(putUserPolicyInput)
	if err != nil {
		return err
	}
	return nil
}
