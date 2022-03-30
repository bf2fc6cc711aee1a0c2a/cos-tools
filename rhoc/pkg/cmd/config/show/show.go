package show

import (
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewConfigShowCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show",
		Long:  "show",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Config.Load()
			if err != nil {
				return err
			}

			rows := []struct {
				Key string `json:"key" header:"Key"`
				Val string `json:"val" header:"Val"`
			}{
				{Key: "APIUrl", Val: c.APIUrl},
				{Key: "AuthURL", Val: c.AuthURL},
				{Key: "MasAuthURL", Val: c.MasAuthURL},
			}

			dump.Table(f.IOStreams.Out, rows)
			f.Logger.Info("")

			return nil
		},
	}

	return cmd
}
