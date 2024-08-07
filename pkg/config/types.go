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
	// GithubRepos represents the Github repository configuration
	GithubRepos []types.GithubRepo `yaml:"githubRepos"`
	// K8sClusters represents the Kubernetes clusters configuration
	K8sClusters []types.K8sCluster `yaml:"k8sClusters"`
	// OnePasswords represents the 1Password configuration
	OnePasswords []types.OnePassword `yaml:"onePasswords"`
	// Web represents the web configuration
	Web types.Web `yaml:"web"`
}
