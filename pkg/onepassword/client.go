package onepassword

import (
	"context"

	onepasswordsdk "github.com/1password/onepassword-sdk-go"
	"github.com/rusik69/iamrotator/pkg/types"
)

// CreateClient creates a new 1Password client
func CreateClient(cfg types.OnePassword) {
	client, err := onepasswordsdk.NewClient(
		context.TODO(),
		onepasswordsdk.WithServiceAccountToken(cfg.Token),
		onepasswordsdk.WithIntegrationInfo("iamrotator", "v0.0.1"),
	)
}
