package token

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewConfigTokenCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "token",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := service.API(&service.Config{
				F: f,
			})
			if err != nil {
				return err
			}

			mas, err := cmd.Flags().GetBool("mas")
			if err != nil {
				return err
			}

			if mas {
				_, _ = f.IOStreams.Out.Write([]byte(a.GetConfig().MasAccessToken))
			} else {
				_, _ = f.IOStreams.Out.Write([]byte(a.GetConfig().AccessToken))
			}

			return nil
		},
	}

	cmd.Flags().BoolP("mas", "m", false, "mas")

	return cmd
}
