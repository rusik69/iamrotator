package types

import "time"

// AWSConfig represents the AWSConfig configuration
type AWSConfig struct {
	// Name represents the AWS org name
	Name string `yaml:"name"`
	// AccountID represents the AWS account ID
	AccountID string `yaml:"accountID"`
	// Region represents the AWS region
	Region string `yaml:"region"`
	// AccessKeyID represents the AWS access key ID
	AccessKeyID string `yaml:"accessKeyID"`
	// SecretAccessKey represents the AWS secret access key
	SecretAccessKey string `yaml:"secretAccessKey"`
	// SessionToken represents the AWS session token
	SessionToken string `yaml:"sessionToken"`
	// RoleName represents the AWS role name
	RoleName string `yaml:"roleName"`
	// IamUserName represents the AWS IAM user name
	IamUserName string `yaml:"iamUserName"`
	// StackSetName represents the AWS stack set name
	StackSetName string `yaml:"stackSetName"`
	// OUID represents the AWS organization unit ID
	OUID string `yaml:"ouID"`
}

// AWSAccessKey represents the AWS access key
type AWSAccessKey struct {
	// AccessKeyID represents the AWS access key ID
	AccessKeyID string `yaml:"accessKeyID"`
	// SecretAccessKey represents the AWS secret access key
	SecretAccessKey string `yaml:"secretAccessKey"`
	// UserName represents the AWS user name
	UserName string `yaml:"userName"`
	// AccountID represents the AWS account ID
	AccountID string `yaml:"accountID"`
	// AccountName represents the AWS account name
	AccountName string `yaml:"accountName"`
	// CreateDate represents the AWS access key creation date
	CreateDate time.Time `yaml:"createDate"`
	// LastUsedDate represents the AWS access key last used date
	LastUsedDate time.Time `yaml:"lastUsedDate"`
	// Status represents the AWS access key status
	Status string `yaml:"status"`
}

// GerritRepo represents the Gerrit repository configuration
type GerritRepo struct {
	// Name represents the Gerrit repository name
	Name string `yaml:"name"`
	// URL represents the Gerrit repository URL
	URL string `yaml:"url"`
}

// GithubOrg represents the Github organization configuration
type GithubOrg struct {
	// Name represents the Github organization name
	Name string `yaml:"name"`
	// Token represents the Github organization token
	Token string `yaml:"token"`
}

// GithubRepo represents the Github repository configuration
type GithubRepo struct {
	// Name represents the Github repository name
	Name string `yaml:"name"`
	// Token represents the Github repository token
	Token string `yaml:"token"`
	// Owner represents the Github repository owner
	Owner string `yaml:"owner"`
}

// K8sCluster represents the Kubernetes cluster configuration
type K8sCluster struct {
	// Name represents the Kubernetes cluster name
	Name string `yaml:"name"`
	// Kubeconfig represents the Kubernetes cluster kubeconfig
	Kubeconfig string `yaml:"kubeconfig"`
}

// 1Password represents the 1Password configuration
type OnePassword struct {
	// Name represents the 1Password name
	Name string `yaml:"name"`
	// Token represents the 1Password token
	Token string `yaml:"token"`
	// Vault represents the 1Password vault
	Vault string `yaml:"vault"`
}

// Web represents the web configuration
type Web struct {
	ListenAddr        string   `yaml:"listenAddr"`
	SSOClientID       string   `yaml:"ssoClientID"`
	SSOClientSecret   string   `yaml:"ssoClientSecret"`
	SSOCallbackURL    string   `yaml:"ssoCallbackURL"`
	SSOStateString    string   `yaml:"ssoStateString"`
	SSOAllowedEmails  []string `yaml:"ssoAllowedEmails"`
	SSOCookieStoreKey string   `yaml:"ssoCookieStoreKey"`
}
