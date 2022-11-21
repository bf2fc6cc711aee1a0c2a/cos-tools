package patch

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
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
			return opts.ready()
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddYes(cmd, &opts.skipConfirm)
	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}

func (opts options) ready() error {
	if !opts.skipConfirm {
		confirm, promptErr := cmdutil.PromptConfirm("Are you sure you want to start the connector with id '%s'?", opts.id)
		if promptErr != nil {
			return promptErr
		}
		if !confirm {
			opts.f.Logger.Debug("User has chosen to not start connector")
			return nil
		}
	}

	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	reqBody := map[string]interface{}{
		"desired_state": admin.CONNECTORSTATE_READY,
	}

	res, resp, err := c.Clusters().PatchConnector(opts.f.Context, opts.id).Body(reqBody).Execute()
	if err != nil {
		return response.Error(err, resp)
	}

	if resp != nil && resp.StatusCode > 299 {
		return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, res)
	}

	return nil
}
