package update

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName = "stop"
)

type options struct {
	outputFormat string
	id           string

	f *factory.Factory
}

func NewStopCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:  CommandName,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.ValidateOutputs(cmd); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return stop(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}

func stop(opts *options) error {
	c, err := service.NewPublicClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	r := c.ConnectorsApi().PatchConnector(opts.f.Context, opts.id).Body(map[string]interface{}{"desired_state": "stopped"})
	result, httpRes, err := r.Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return response.Error(err, httpRes)
	}

	if httpRes != nil && httpRes.StatusCode == 204 {
		return nil
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, result)
}
