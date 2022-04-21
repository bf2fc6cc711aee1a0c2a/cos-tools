package token

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewConfigTokenCommand(f *factory.Factory) *cobra.Command {
	var mas bool

	cmd := &cobra.Command{
		Use:   "token",
		Short: "token",
		Long:  "token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := service.API(&service.Config{
				F: f,
			})
			if err != nil {
				return err
			}

			if mas {
				f.Logger.Info(a.GetConfig().MasAccessToken)
			} else {
				f.Logger.Info(a.GetConfig().AccessToken)
			}

			return nil
		},
	}

	flags := flagutil.NewFlagSet(cmd, f.Localizer)
	flags.BoolVarP(&mas, "mas", "m", false, "mas")

	return cmd
}
