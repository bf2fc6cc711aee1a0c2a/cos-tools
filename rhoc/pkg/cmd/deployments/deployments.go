package deployments

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/get"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/list"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/snapshot"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/updateChannel"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/validateSnapshot"
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
		updateChannel.NewUpdateChannelCommand(f),
		snapshot.NewSnapshotCommand(f),
		validateSnapshot.NewValidateSnapshotCommand(f))

	return cmd
}
