package describe

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName  = "describe"
	CommandAlias = "get"
)

type options struct {
	id           string
	clusterID    string
	outputFormat string

	f *factory.Factory
}

func NewDescribeCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     CommandName,
		Aliases: []string{CommandAlias},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.ValidateOutputs(cmd); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddID(cmd, &opts.id).Required()
	cmdutil.AddClusterID(cmd, &opts.clusterID).Required()

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	result, httpRes, err := c.Clusters().GetConnectorDeployment(opts.f.Context, opts.clusterID, opts.id).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return response.Error(err, httpRes)
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, result)
}
