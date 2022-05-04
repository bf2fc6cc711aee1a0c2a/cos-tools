package whoami

import (
	"fmt"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/core/auth/token"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewWhoAmICommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "whoami",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			a, err := service.API(&service.Config{
				F: f,
			})
			if err != nil {
				return err
			}

			userName, ok := token.GetUsername(a.GetConfig().MasAccessToken)
			if !ok {
				userName = "unknown"
			}

			if ok {
				fmt.Fprintf(f.IOStreams.Out, "%s@%s\n", userName, a.GetConfig().ApiURL.String())
			}

			return nil
		},
	}

	return cmd
}
