package main

import (
	"github.com/rusik69/iamrotator/pkg/aws"
	"github.com/rusik69/iamrotator/pkg/config"
	"github.com/rusik69/iamrotator/pkg/github"

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
			keys, err := aws.ListAccessKeys(awsSess, cfg.AWS)
			if err != nil {
				logrus.Panic(err)
			}
			for _, key := range keys {
				logrus.Info("Account ID:", key.AccountID, "User Name:", key.UserName, "Access Key ID:", key.AccessKeyID)
			}
		},
	}

	var listOrgSecretsCmd = &cobra.Command{
		Use:   "listgithuborgsecrets",
		Short: "List github organization secrets",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(configPath)
			if err != nil {
				logrus.Panic(err)
			}
			for _, org := range cfg.GithubOrgs {
				orgSess, err := github.CreateOrgClient(org)
				if err != nil {
					logrus.Panic(err)
				}
				logrus.Info("Listing organization secrets")
				secrets, err := github.ListOrgSecrets(orgSess, org)
				if err != nil {
					logrus.Panic(err)
				}
				for _, secret := range secrets {
					logrus.Info("Secret Name:", secret)
				}
			}
		},
	}

	var listRepoSecretsCmd = &cobra.Command{
		Use:   "listgithubreposecrets",
		Short: "List github repository secrets",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(configPath)
			if err != nil {
				logrus.Panic(err)
			}
			for _, repo := range cfg.GithubRepos {
				repoSess, err := github.CreateRepoClient(repo)
				if err != nil {
					logrus.Panic(err)
				}
				logrus.Info("Listing repository secrets")
				secrets, err := github.ListRepoSecrets(repoSess, repo)
				if err != nil {
					logrus.Panic(err)
				}
				for _, secret := range secrets {
					logrus.Info("Secret Name:", secret)
				}
			}
		},
	}

	rootCmd.AddCommand(createUserCmd, createStackSetCmd, removeUserCmd, removeStackSetCmd, listKeysCmd, listOrgSecretsCmd, listRepoSecretsCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Panic(err)
	}
}
