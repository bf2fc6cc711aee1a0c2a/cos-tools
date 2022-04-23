package namespaces

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/create"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/delete"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/describe"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/list"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/commands"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewNamespacesCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespaces",
		Aliases: []string{"ns"},
		Short:   "namespaces",
		Long:    "namespaces",
		Args:    cobra.MinimumNArgs(1),
	}

	commands.Bind(
		cmd,
		list.NewListCommand(f),
		create.NewCreateCommand(f),
		delete.NewDeletesCommand(f),
		describe.NewDescribeCommand(f))

	return cmd
}
