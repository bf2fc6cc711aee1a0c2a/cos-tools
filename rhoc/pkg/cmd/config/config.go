package config

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/config/context"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/config/show"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/config/token"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

type itemRow struct {
	Key string `json:"key" header:"Key"`
	Val string `json:"val" header:"Val"`
}

func NewConfigCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(show.NewConfigShowCommand(f))
	cmd.AddCommand(token.NewConfigTokenCommand(f))
	cmd.AddCommand(context.NewContextCommand(f))

	return cmd
}
