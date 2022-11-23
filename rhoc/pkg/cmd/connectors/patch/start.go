package patch

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewStartCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use: "start",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.ValidateOutputs(cmd); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.skipConfirm {
				confirm, err := cmdutil.PromptConfirm("Are you sure you want to start the connector with id '%s'?", opts.id)
				if err != nil {
					return err
				}
				if !confirm {
					opts.f.Logger.Debugf("User has chosen to not start connector")
					return nil
				}
			}

			return execute(opts, map[string]interface{}{
				"desired_state": admin.CONNECTORDESIREDSTATE_READY,
			})
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddYes(cmd, &opts.skipConfirm)
	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}
