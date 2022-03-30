package token

import (
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
			c, err := f.Config.Load()
			if err != nil {
				return err
			}

			if mas {
				f.Logger.Info(c.MasAccessToken)
			} else {
				f.Logger.Info(c.AccessToken)
			}

			return nil
		},
	}

	flags := flagutil.NewFlagSet(cmd, f.Localizer)
	flags.BoolVarP(&mas, "mas", "m", false, "mas")

	return cmd
}
