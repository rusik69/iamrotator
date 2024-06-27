package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/rusik69/iamrotator/pkg/types"
	"github.com/sirupsen/logrus"
)

// ListStackSets lists all stack sets
func ListStackSets(sess *session.Session) ([]string, error) {
	svc := cloudformation.New(sess)
	input := &cloudformation.ListStackSetsInput{
		Status: aws.String("ACTIVE"),
	}
	result, err := svc.ListStackSets(input)
	if err != nil {
		return nil, err
	}
	var stackSetNames []string
	for _, stackSet := range result.Summaries {
		stackSetNames = append(stackSetNames, *stackSet.StackSetName)
	}
	return stackSetNames, nil
}

// CreateRoleStackSet creates a new stack set for the role
func CreateRoleStackSet(sess *session.Session, cfg types.AWS) error {
	svc := cloudformation.New(sess)
	principalArn := "arn:aws:iam::" + cfg.AccountID + ":user/" + cfg.IamUserName
	template := `{
        "Resources": {
            "IAMAccessRole": {
                "Type": "AWS::IAM::Role",
                "Properties": {
                    "RoleName": "IAMAccessRole",
                    "AssumeRolePolicyDocument": {
                        "Version": "2012-10-17",
                        "Statement": [
                            {
                                "Effect": "Allow",
                                "Principal": {
                                    "AWS": "` + principalArn + `"
                                },
                                "Action": "sts:AssumeRole"
                            }
                        ]
                    },
                    "ManagedPolicyArns": [
                        "arn:aws:iam::aws:policy/IAMFullAccess"
                    ]
                }
            }
        }
    }`
	input := &cloudformation.CreateStackSetInput{
		StackSetName: aws.String("iamrotator"),
		TemplateBody: aws.String(template),
	}
	fmt.Println("Creating stack set iamrotator")
	_, err := svc.CreateStackSet(input)
	if err != nil {
		return err
	}
	orgs := organizations.New(sess)
	// Slice to hold all account IDs
	var accountIDs []*string

	// Handle pagination
	var nextToken *string
	for {
		accountsRes, err := orgs.ListAccounts(&organizations.ListAccountsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return err
		}
		for _, account := range accountsRes.Accounts {
			accountIDs = append(accountIDs, account.Id)
		}
		if accountsRes.NextToken == nil {
			break
		}
		nextToken = accountsRes.NextToken
	}
	createStackInstancesInput := &cloudformation.CreateStackInstancesInput{
		StackSetName: aws.String("iamrotator"),
		Accounts:     accountIDs,
		Regions:      []*string{aws.String("us-east-1")},
	}
	fmt.Println("Creating stack instances")
	_, err = svc.CreateStackInstances(createStackInstancesInput)
	if err != nil {
		return err
	}
	for {
		describeStackSetOperationInput := &cloudformation.DescribeStackSetOperationInput{
			StackSetName: aws.String("iamrotator"),
			OperationId:  aws.String("iamrotator"),
		}
		result, err := svc.DescribeStackSetOperation(describeStackSetOperationInput)
		if err != nil {
			logrus.Errorf("Failed to describe stack set operation: %v", err)
		}
		if *result.StackSetOperation.Status == "SUCCEEDED" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}

// CheckOrCreateStackSet checks if the stack set exists and creates it if it doesn't
func CheckOrCreateStackSet(sess *session.Session, cfg types.AWS) error {
	stackSets, err := ListStackSets(sess)
	if err != nil {
		return err
	}
	stackSetFound := false
	for _, stackSet := range stackSets {
		if stackSet == "iamrotator" {
			stackSetFound = true
			break
		}
	}
	if stackSetFound {
		fmt.Println("Stack set iamrotator found")
		return nil
	}
	logrus.Info("Stack set iamrotator not found")
	return CreateRoleStackSet(sess, cfg)
}

// EmptyStackSet empties the stack set
func EmptyStackSet(sess *session.Session, stackSetName, region string) error {
	svc := cloudformation.New(sess)
	input := &cloudformation.ListStackInstancesInput{
		StackSetName:        aws.String(stackSetName),
		StackInstanceRegion: aws.String(region),
	}
	for {
		failed := false
		fmt.Println("Listing stack instances")
		result, err := svc.ListStackInstances(input)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", result)
		for _, instance := range result.Summaries {
			fmt.Println("Deleting stack instance", *instance.Account, *instance.Region)
			deleteInput := &cloudformation.DeleteStackInstancesInput{
				StackSetName: aws.String(stackSetName),
				Accounts:     []*string{instance.Account},
				Regions:      []*string{instance.Region},
				RetainStacks: aws.Bool(false),
			}
			_, err := svc.DeleteStackInstances(deleteInput)
			if err != nil {
				logrus.Error(err)
				failed = true
			}
		}
		if !failed {
			break
		} else {
			logrus.Info("Retrying in 10 seconds")
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}

// RemoveStackSet removes the stack set
func RemoveStackSet(sess *session.Session, stackSetName string) error {
	svc := cloudformation.New(sess)
	for {
		failed := false
		input := &cloudformation.DeleteStackSetInput{
			StackSetName: aws.String(stackSetName),
		}
		_, err := svc.DeleteStackSet(input)
		if err != nil {
			return err
		}
		if !failed {
			break
		} else {
			logrus.Info("Retrying in 10 seconds")
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}
