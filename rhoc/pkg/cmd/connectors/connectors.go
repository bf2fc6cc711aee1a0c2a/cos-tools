package connectors

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors/delete"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors/get"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors/list"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors/types"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewConnectorsCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connectors",
		Aliases: []string{"ctr"},
		Args:    cobra.MinimumNArgs(1),
	}

	cmdutil.Bind(
		cmd,
		list.NewListCommand(f),
		delete.NewDeletesCommand(f),
		get.NewGetCommand(f),
		types.NeConnectorTypesCommand(f))

	return cmd
}
