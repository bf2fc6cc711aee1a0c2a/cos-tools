package service

import (
	"context"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/public"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/auth"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context) (*public.APIClient, error) {
	var token string
	var err error

	token = viper.GetString("api-token")
	if token == "" {
		token, err = auth.Token()
		if err != nil {
			return nil, err
		}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	client := public.NewAPIClient(&public.Configuration{
		BasePath:   viper.GetString("api-url"),
		HTTPClient: oauth2.NewClient(ctx, ts),
	})

	return client, nil
}
