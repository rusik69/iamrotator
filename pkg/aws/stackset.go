package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/rusik69/iamrotator/pkg/types"
	"github.com/sirupsen/logrus"
)

// ListStackSets lists all stack sets
func ListStackSets(sess aws.Config) ([]string, error) {
	svc := cloudformation.NewFromConfig(sess)
	input := &cloudformation.ListStackSetsInput{
		Status: "ACTIVE",
		CallAs: cftypes.CallAsDelegatedAdmin,
	}
	result, err := svc.ListStackSets(context.TODO(), input)
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
func CreateRoleStackSet(sess aws.Config, cfg types.AWS) error {
	svc := cloudformation.NewFromConfig(sess)
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
		StackSetName: aws.String(cfg.StackSetName),
		TemplateBody: aws.String(template),
		AutoDeployment: &cftypes.AutoDeployment{
			Enabled:                      aws.Bool(true),
			RetainStacksOnAccountRemoval: aws.Bool(false),
		},
		PermissionModel: cftypes.PermissionModelsServiceManaged,
		CallAs:          cftypes.CallAsDelegatedAdmin,
	}
	ctx := context.Background()
	logrus.Info("Creating stack set", cfg.StackSetName)
	_, err := svc.CreateStackSet(ctx, input)
	if err != nil {
		return err
	}
	orgs := organizations.NewFromConfig(sess)
	rootID := ""
	rootInput := &organizations.ListRootsInput{}
	rootResult, err := orgs.ListRoots(ctx, rootInput)
	if err != nil {
		return err
	}
	for _, root := range rootResult.Roots {
		rootID = *root.Id
		break
	}
	logrus.Info("Root ID:", rootID)
	createStackInstancesInput := &cloudformation.CreateStackInstancesInput{
		StackSetName: aws.String(cfg.StackSetName),
		DeploymentTargets: &cftypes.DeploymentTargets{
			OrganizationalUnitIds: []string{rootID},
		},
		Regions: []string{cfg.Region},
		CallAs:  cftypes.CallAsDelegatedAdmin,
	}
	logrus.Info("Creating stack set instances")
	out, err := svc.CreateStackInstances(ctx, createStackInstancesInput)
	if err != nil {
		return err
	}
	for {
		describeStackSetOperationInput := &cloudformation.DescribeStackSetOperationInput{
			StackSetName: aws.String(cfg.StackSetName),
			OperationId:  out.OperationId,
			CallAs:       cftypes.CallAsDelegatedAdmin,
		}
		result, err := svc.DescribeStackSetOperation(ctx, describeStackSetOperationInput)
		if err != nil {
			logrus.Errorf("Failed to describe stack set operation: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if result.StackSetOperation.Status == "SUCCEEDED" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}

// CheckOrCreateStackSet checks if the stack set exists and creates it if it doesn't
func CheckOrCreateStackSet(sess aws.Config, cfg types.AWS) error {
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
		logrus.Println("Stack set iamrotator found")
		return nil
	}
	logrus.Info("Stack set iamrotator not found")
	return CreateRoleStackSet(sess, cfg)
}

// EmptyStackSet empties the stack set
func EmptyStackSet(sess aws.Config, stackSetName, region string) error {
	svc := cloudformation.NewFromConfig(sess)
	logrus.Info("Emptying stack set", stackSetName, region)
	deleteInput := cloudformation.DeleteStackInstancesInput{
		StackSetName: aws.String(stackSetName),
		Regions:      []string{region},
		RetainStacks: aws.Bool(false),
		CallAs:       cftypes.CallAsDelegatedAdmin,
	}
	_, err := svc.DeleteStackInstances(context.TODO(), &deleteInput)
	if err != nil {
		return err
	}
	return nil
}

// RemoveStackSet removes the stack set
func RemoveStackSet(sess aws.Config, stackSetName string) error {
	svc := cloudformation.NewFromConfig(sess)
	for {
		failed := false
		input := &cloudformation.DeleteStackSetInput{
			StackSetName: aws.String(stackSetName),
			CallAs:       cftypes.CallAsDelegatedAdmin,
		}
		_, err := svc.DeleteStackSet(context.TODO(), input)
		if err != nil {
			failed = true
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
