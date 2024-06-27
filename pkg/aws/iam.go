package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
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
	logrus.Info("Creating IAM user", userName)
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

// CreateAccessKeys creates access keys for the IAM user
func CreateAccessKeys(sess *session.Session, userName string) (string, string, error) {
	svc := iam.New(sess)
	input := &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	}
	result, err := svc.CreateAccessKey(input)
	if err != nil {
		return "", "", err
	}
	return *result.AccessKey.AccessKeyId, *result.AccessKey.SecretAccessKey, err
}

// CheckOrCreateIamUser checks if the IAM user exists and creates it if it doesn't
func CheckOrCreateIamUser(sess *session.Session, userName string) (string, string, error) {
	users, err := ListIamUsers(sess)
	if err != nil {
		return "", "", err
	}
	userFound := false
	for _, user := range users {
		if user == userName {
			userFound = true
			break
		}
	}
	if userFound {
		logrus.Info("IAM user", userName, "found")
		return "", "", err
	}
	logrus.Info("IAM user", userName, "not found")
	err = CreateIamUser(sess, userName)
	if err != nil {
		return "", "", err
	}
	return CreateAccessKeys(sess, userName)
}

// RemoveIamUser removes the IAM user
func RemoveIamUser(sess *session.Session, userName string) error {
	svc := iam.New(sess)
	listPoliciesInput := &iam.ListUserPoliciesInput{
		UserName: aws.String(userName),
	}
	policyNames, err := svc.ListUserPolicies(listPoliciesInput)
	if err != nil {
		return err
	}
	for _, policyName := range policyNames.PolicyNames {
		logrus.Info("Detaching policy", *policyName)
		detachPolicyInput := &iam.DeleteUserPolicyInput{
			UserName:   aws.String(userName),
			PolicyName: policyName,
		}
		_, err := svc.DeleteUserPolicy(detachPolicyInput)
		if err != nil {
			return err
		}
	}
	listAccessKeysInput := &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}
	accessKeys, err := svc.ListAccessKeys(listAccessKeysInput)
	if err != nil {
		return err
	}
	for _, accessKey := range accessKeys.AccessKeyMetadata {
		deleteAccessKeyInput := &iam.DeleteAccessKeyInput{
			AccessKeyId: accessKey.AccessKeyId,
			UserName:    aws.String(userName),
		}
		logrus.Info("Deleting access key", *accessKey.AccessKeyId)
		_, err := svc.DeleteAccessKey(deleteAccessKeyInput)
		if err != nil {
			return err
		}
	}
	input := &iam.DeleteUserInput{
		UserName: &userName,
	}
	_, err = svc.DeleteUser(input)
	if err != nil {
		return err
	}
	return nil
}
