package deployments

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/get"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/list"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/update"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewDeploymentsCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployments",
		Aliases: []string{"cd"},
		Args:    cobra.MinimumNArgs(1),
	}

	cmdutil.Bind(
		cmd,
		list.NewListCommand(f),
		get.NewGetCommand(f),
		update.NewUpdateCommand(f))

	return cmd
}
