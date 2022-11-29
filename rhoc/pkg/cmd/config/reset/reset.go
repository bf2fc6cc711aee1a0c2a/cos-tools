package reset

import (
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewConfigResetCommand(f *factory.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:  "reset",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Config.Remove(); err != nil {
				return err
			}

			if err := f.ServiceContext.Remove(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
