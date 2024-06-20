package types

// AWS represents the AWS configuration
type AWS struct {
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
	// User represents the Github organization user
	User string `yaml:"user"`
}

// K8sCluster represents the Kubernetes cluster configuration
type K8sCluster struct {
	// Name represents the Kubernetes cluster name
	Name string `yaml:"name"`
	// Kubeconfig represents the Kubernetes cluster kubeconfig
	Kubeconfig string `yaml:"kubeconfig"`
}
