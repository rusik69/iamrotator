package github

import (
	"context"

	"github.com/google/go-github/v63/github"
	"github.com/rusik69/iamrotator/pkg/types"
)

// ListRepoSecrets lists the Github repository secrets
func ListRepoSecrets(client *github.Client, repo types.GithubRepo) ([]string, error) {
	ctx := context.Background()
	secrets, _, err := client.Actions.ListRepoSecrets(ctx, repo.Owner, repo.Name, nil)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, secret := range secrets.Secrets {
		res = append(res, secret.Name)
	}
	return res, nil
}
