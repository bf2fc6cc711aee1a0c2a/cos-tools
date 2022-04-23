package delete

import (
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"net/http"
)

const (
	CommandName  = "delete"
	CommandAlias = "rm"
)

type options struct {
	outputFormat string
	id           string
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
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, cmdutil.ValidOutputs()...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, cmdutil.ValidOutputs()...)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddYes(cmd, &opts.skipConfirm)
	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}

func run(opts *options) error {
	if !opts.skipConfirm {
		confirm, promptErr := cmdutil.PromptConfirm("Are you sure you want to delete the connector namespace with id '%s'?", opts.id)
		if promptErr != nil {
			return promptErr
		}
		if !confirm {
			opts.f.Logger.Debug("User has chosen to not delete connector namespace")
			return nil
		}
	}

	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	e := c.ConnectorClustersAdminApi.DeleteConnectorNamespace(opts.f.Context, opts.id)

	response, httpRes, err := e.Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		if httpRes != nil && httpRes.StatusCode == http.StatusInternalServerError {
			e, _ := service.ReadError(httpRes)
			if e.Reason != "" {
				err = fmt.Errorf("%s: [%w]", err.Error(), errors.New(e.Reason))
			}
		}
		return err
	}
	if httpRes != nil && httpRes.StatusCode == 204 {
		return nil
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, response)
}
