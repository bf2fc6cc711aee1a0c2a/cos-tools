package delete

import (
	"strconv"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName  = "delete"
	CommandAlias = "rm"
)

type options struct {
	outputFormat string
	id           string
	force        bool
	skipConfirm  bool

	f *factory.Factory
}

func NewDeletesCommand(f *factory.Factory) *cobra.Command {
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
	cmdutil.AddYes(cmd, &opts.skipConfirm)
	cmdutil.AddForce(cmd, &opts.force)
	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}

func run(opts *options) error {
	if !opts.skipConfirm {
		confirm, promptErr := cmdutil.PromptConfirm("Are you sure you want to delete the connector with id '%s'?", opts.id)
		if promptErr != nil {
			return promptErr
		}
		if !confirm {
			opts.f.Logger.Debug("User has chosen to not delete connector")
			return nil
		}
	}

	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	e := c.Clusters().DeleteConnector(opts.f.Context, opts.id)
	e = e.Force(strconv.FormatBool(opts.force))

	result, httpRes, err := e.Execute()
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
