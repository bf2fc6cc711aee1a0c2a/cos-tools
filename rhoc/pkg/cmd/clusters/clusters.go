package clusters

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/clusters/get"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/clusters/list"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/clusters/tree"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NeClustersCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clusters",
		Aliases: []string{"cc"},
		Args:    cobra.MinimumNArgs(1),
	}

	return cmdutil.Bind(
		cmd,
		list.NewListCommand(f),
		tree.NewTreeCommand(f),
		get.NewGetCommand(f))
}
