package namespaces

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/commands"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/create"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/delete"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces/list"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewNamespacesCommand(f *factory.Factory) *cobra.Command {
	ns := &cobra.Command{
		Use:     "namespaces",
		Aliases: []string{"ns"},
		Short:   "ns",
		Long:    "ns",
		Args:    cobra.MinimumNArgs(1),
	}

	commands.Bind(
		ns,
		list.NewListCommand(f),
		create.NewCreateCommand(f),
		delete.NewDeletesCommand(f))

	return ns
}
