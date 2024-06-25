package main

import (
	"fmt"
	"os"

	"github.com/rusik69/iamrotator/pkg/args"
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
)

func main() {
	arg := args.Parse()
	cfg, err := config.Load(arg.ConfigPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded configuration from %s\n", arg.ConfigPath)
	fmt.Printf("Config: %v\n", cfg)
	awsSess, err := aws.CreateSession(cfg.AWS)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Checking AWS account %s\n", cfg.AWS.AccountID)
	if arg.Action == "createuser" {
		newKeyID, NewKeySecret, err := aws.CheckOrCreateIamUser(awsSess, cfg.AWS.IamUserName)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Access Key ID: %s\n", newKeyID)
		fmt.Printf("Secret Access Key: %s\n", NewKeySecret)
	} else if arg.Action == "createstackset" {
		fmt.Println("Checking or creating stack set")
		err = aws.CheckOrCreateStackSet(awsSess, cfg.AWS)
		if err != nil {
			panic(err)
		}
	} else if arg.Action == "removeuser" {
		fmt.Println("Removing IAM user")
		err = aws.RemoveIamUser(awsSess, cfg.AWS.IamUserName)
		if err != nil {
			panic(err)
		}
	} else if arg.Action == "removestackset" {
		fmt.Println("Removing stack set")
		err = aws.EmptyStackSet(awsSess, "iamrotator", cfg.AWS.Region)
		if err != nil {
			panic(err)
		}
		err = aws.RemoveStackSet(awsSess, "iamrotator")
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Usage: iamrotator <action> <configpath>")
		os.Exit(1)
	}
}
