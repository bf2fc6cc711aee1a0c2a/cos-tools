package version

import (
	"fmt"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewVersionCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: false,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(f.IOStreams.Out, build.Version)
			return err
		},
	}

	return cmd
}
