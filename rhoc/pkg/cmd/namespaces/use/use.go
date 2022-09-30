package use

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/contextutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName = "use"
)

type options struct {
	id string

	f *factory.Factory
}

func NewUseCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:  CommandName,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}

	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}

func run(opts *options) error {

	svcContext, err := opts.f.ServiceContext.Load()
	if err != nil {
		return err
	}

	currCtx, err := contextutil.GetCurrentContext(svcContext, opts.f.Localizer)
	if err != nil {
		return err
	}

	currCtx.NamespaceID = opts.id
	svcContext.Contexts[svcContext.CurrentContext] = *currCtx

	if err := opts.f.ServiceContext.Save(svcContext); err != nil {
		return err
	}

	return nil
}
