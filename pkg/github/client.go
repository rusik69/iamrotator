package github

import (
	"context"

	"github.com/google/go-github/v63/github"
	"github.com/rusik69/iamrotator/pkg/types"
	"golang.org/x/oauth2"
)

// CreateOrgClient creates a new Github org client
func CreateOrgClient(cfg types.GithubOrg) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, nil
}

// CreateRepoClient creates a new Github repo client
func CreateRepoClient(cfg types.GithubRepo) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, nil
}
