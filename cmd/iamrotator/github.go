package main

import (
	"github.com/rusik69/iamrotator/pkg/config"
	"github.com/rusik69/iamrotator/pkg/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
