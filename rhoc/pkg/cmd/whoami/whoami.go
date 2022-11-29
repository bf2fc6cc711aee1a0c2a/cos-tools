package whoami

import (
	"fmt"
	"strings"

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

			if a.GetConfig().AccessToken == "" {
				return nil
			}

			accessTkn, err := token.Parse(a.GetConfig().AccessToken)
			if err != nil {
				return err
			}

			tknClaims, _ := token.MapClaims(accessTkn)

			realms, ok := tknClaims["realm_access"].(map[string]interface{})
			if !ok {
				return nil
			}
			roles, ok := realms["roles"].([]interface{})
			if !ok {
				return nil
			}

			_, _ = fmt.Fprintln(f.IOStreams.Out, "roles:")
			for i := range roles {
				role := roles[i].(string)

				if strings.HasPrefix(role, "cos-fleet-manager-") {
					_, _ = fmt.Fprintf(f.IOStreams.Out, "  %s\n", role)
				}
			}

			return nil
		},
	}

	return cmd
}
