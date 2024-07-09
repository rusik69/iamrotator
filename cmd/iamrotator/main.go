package main

import (
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "iamrotator [createuser, removeuser, createstackset, removestackset] -c <configpath>",
		Short: "IAM Rotator is a tool to manage AWS IAM users and stack sets",
	}

	var configPath string

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to the configuration file")
	rootCmd.MarkPersistentFlagRequired("config")

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
			err, keys = aws.ListAccessKeys(awsSess, cfg.AWS)
			if err != nil {
				logrus.Panic(err)
			}
		},
	}

	rootCmd.AddCommand(createUserCmd, createStackSetCmd, removeUserCmd, removeStackSetCmd, listKeysCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Panic(err)
	}
}
