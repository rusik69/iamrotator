package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "iamrotator [createuser, removeuser, createstackset, removestackset] -c <configpath>",
		Short: "IAM Rotator is a tool to manage AWS IAM users and stack sets",
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to the configuration file")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(createUserCmd, createStackSetCmd, removeUserCmd, removeStackSetCmd, listKeysCmd, CreateAccessKeyCmd, listOrgSecretsCmd, listRepoSecretsCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Panic(err)
	}
}
