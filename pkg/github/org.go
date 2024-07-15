package github

import (
	"context"

	"github.com/google/go-github/v63/github"
	"github.com/rusik69/iamrotator/pkg/types"
)

// ListOrgSecrets lists the Github organization secrets
func ListOrgSecrets(client *github.Client, org types.GithubOrg) ([]string, error) {
	ctx := context.Background()
	secrets, _, err := client.Actions.ListOrgSecrets(ctx, org.Name, nil)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, secret := range secrets.Secrets {
		res = append(res, secret.Name)
	}
	return res, nil
}
