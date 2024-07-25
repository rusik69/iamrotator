package main

import (
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
	"github.com/rusik69/iamrotator/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath, userName string

var createUserCmd = &cobra.Command{
	Use:   "createuser",
	Short: "Create a new IAM user",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logrus.Panic(err)
		}
		awsSess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Checking AWS account ", cfg.AWS.AccountID)
		newKeyID, NewKeySecret, err := aws.CheckOrCreateIamUser(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Access Key ID:", newKeyID)
		logrus.Info("Secret Access Key:", NewKeySecret)
	},
}

var removeUserCmd = &cobra.Command{
	Use:   "removeuser",
	Short: "Remove an IAM user",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logrus.Panic(err)
		}
		awsSess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Removing IAM user")
		err = aws.RemoveIamUser(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("IAM user removed")
	},
}

var createStackSetCmd = &cobra.Command{
	Use:   "createstackset",
	Short: "Create a new stack set",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logrus.Panic(err)
		}
		awsSess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Checking or creating stack set")
		err = aws.CheckOrCreateStackSet(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
	},
}

var removeStackSetCmd = &cobra.Command{
	Use:   "removestackset",
	Short: "Remove a stack set",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logrus.Panic(err)
		}
		awsSess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		err = aws.EmptyStackSet(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Removing stack set")
		err = aws.RemoveStackSet(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Stack set removed")
	},
}

var listKeysCmd = &cobra.Command{
	Use:   "listkeys",
	Short: "List access keys",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logrus.Panic(err)
		}
		awsSess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Listing access keys")
		keys, err := aws.ListAccessKeys(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		for _, key := range keys {
			logrus.Info("Account ID: ", key.AccountID, " User Name: ", key.UserName, " Access Key ID: ", key.AccessKeyID, " Create Date: ", key.CreateDate, " Status: ", key.Status)
		}
	},
}

var createAccessKeyCmd = &cobra.Command{
	Use:   "createaccesskey",
	Short: "Create a new access key for the user",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logrus.Panic(err)
		}
		awsSess, err := aws.CreateSession(cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Listing access keys")
		keys, err := aws.ListAccessKeys(awsSess, cfg.AWS)
		if err != nil {
			logrus.Panic(err)
		}
		var foundKeys []types.AWSAccessKey
		for _, key := range keys {
			if key.UserName == userName {
				foundKeys = append(foundKeys, key)
			}
		}
		if len(foundKeys) == 0 {
			logrus.Panic("User not found")
		}
		if len(foundKeys) == 2 {
			logrus.Error(foundKeys[0].AccessKeyID)
			logrus.Error(foundKeys[1].AccessKeyID)
			logrus.Panic("User already has 2 access keys, remove one first")
		}
		logrus.Info("Creating access key")
		_, err = aws.CreateAccessKey(awsSess, cfg.AWS, userName)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Info("Access key created")
	},
}
