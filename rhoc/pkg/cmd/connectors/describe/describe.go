package describe

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

type options struct {
	id           string
	outputFormat string

	f *factory.Factory
}

func NewDescribeCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     "describe",
		Aliases: []string{"get"},
		Short:   "describe",
		Long:    "describe",
		Example: "describe",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, flagutil.ValidOutputFormats...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, flagutil.ValidOutputFormats...)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	flags := flagutil.NewFlagSet(cmd, f.Localizer)
	flags.AddOutput(&opts.outputFormat)
	flags.StringVar(&opts.id, "id", "", "id")

	cmd.MarkFlagRequired("id")

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	result, httpRes, err := c.ConnectorClustersAdminApi.GetConnector(opts.f.Context, opts.id).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return err
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, result)
}
