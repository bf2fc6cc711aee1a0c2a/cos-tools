package types

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors/types/get"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors/types/list"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NeConnectorTypesCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "types",
		Aliases: []string{"t"},
		Args:    cobra.MinimumNArgs(1),
	}

	cmdutil.Bind(
		cmd,
		list.NewListCommand(f),
		get.NewGetCommand(f))

	return cmd
}
