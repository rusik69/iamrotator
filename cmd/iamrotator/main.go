package main

import (
	"fmt"

	"github.com/rusik69/iamrotator/pkg/args"
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
)

func main() {
	args := args.Parse()
	cfg, err := config.Load(args.ConfigPath)
	if err != nil {
		panic(err)
	}
	for _, awsCfg := range cfg.AWS {
		fmt.Printf("Checking AWS account %s\n", awsCfg.AccountID)
		awsSess, err := aws.CreateSession(awsCfg)
		if err != nil {
			panic(err)
		}
		users, err := aws.ListIamUsers(awsSess)
		if err != nil {
			panic(err)
		}
		userFound := false
		for _, user := range users {
			if user == awsCfg.IamUserName {
				userFound = true
				break
			}
		}
		if userFound {
			fmt.Println("IAM User found")
		} else {
			fmt.Println("IAM User not found, creating...")
			err := aws.CreateIamUser(awsSess, awsCfg.IamUserName)
			if err != nil {
				panic(err)
			}
		}
		ssList, err := aws.ListStackSets(awsSess)
		if err != nil {
			panic(err)
		}
		roleStackSetFound := false
		for _, ss := range ssList {
			if ss == "iamrotator" {
				roleStackSetFound = true
				break
			}
		}
		if roleStackSetFound {
			fmt.Println("Role stack set found")
		} else {
			fmt.Println("Role stack set not found, creating...")
			err := aws.CreateRoleStackSet(awsSess, awsCfg)
			if err != nil {
				panic(err)
			}
		}
	}
}
