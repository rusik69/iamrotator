package config

import "github.com/rusik69/iamrotator/pkg/types"

// Config represents the configuration
type Config struct {
	// AWS represents the AWS configuration
	AWS types.AWSConfig `yaml:"aws"`
	// GerritRepos represents the Gerrit repository configuration
	GerritRepos []types.GerritRepo `yaml:"gerritRepos"`
	// GithubOrgs represents the Github organization configuration
	GithubOrgs []types.GithubOrg `yaml:"githubOrgs"`
	// K8sClusters represents the Kubernetes clusters configuration
	K8sClusters []types.K8sCluster `yaml:"k8sClusters"`
}
