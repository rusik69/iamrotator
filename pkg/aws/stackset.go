package aws

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/rusik69/iamrotator/pkg/types"
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
            "RootAccessRole": {
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
	_, err := svc.CreateStackSet(input)
	if err != nil {
		return err
	}
	createStackInstancesInput := &cloudformation.CreateStackInstancesInput{
		StackSetName: aws.String("iamrotator"),
		Accounts:     []*string{aws.String("*")},
		Regions:      []*string{aws.String("us-east-1")},
	}

	_, err = svc.CreateStackInstances(createStackInstancesInput)
	if err != nil {
		log.Fatalf("Failed to create stack instances: %v", err)
	}
	for {
		describeStackSetOperationInput := &cloudformation.DescribeStackSetOperationInput{
			StackSetName: aws.String("iamrotator"),
			OperationId:  aws.String("iamrotator"),
		}

		result, err := svc.DescribeStackSetOperation(describeStackSetOperationInput)
		if err != nil {
			log.Fatalf("Failed to describe stack set operation: %v", err)
		}

		if *result.StackSetOperation.Status == "SUCCEEDED" {
			break
		}

		time.Sleep(10 * time.Second)
	}
	return nil
}
