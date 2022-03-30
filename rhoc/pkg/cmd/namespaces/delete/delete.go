package delete

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
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
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "delete",
		Long:    "delete",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, flagutil.ValidOutputFormats...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, flagutil.ValidOutputFormats...)
			}

			return run(&opts)
		},
	}

	flags := flagutil.NewFlagSet(cmd, f.Localizer)
	flags.AddOutput(&opts.outputFormat)
	flags.AddYes(&opts.skipConfirm)
	flags.StringVar(&opts.id, "id", "", "id")

	cmd.MarkFlagRequired("id")

	return cmd
}

func run(opts *options) error {
	if !opts.skipConfirm {
		confirm, promptErr := promptConfirmDelete(opts)
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

	response, httpRes, err := c.ConnectorClustersAdminApi.DeleteConnectorNamespace(opts.f.Context, opts.id)
	if httpRes != nil {
		defer httpRes.Body.Close()
	}
	if err != nil {
		return err
	}

	if err = dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, response); err != nil {
		return err
	}

	return nil
}

func promptConfirmDelete(opts *options) (bool, error) {
	promptConfirm := survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete the connector namespace with id '%s'?", opts.id),
	}

	var confirmDelete bool
	if err := survey.AskOne(&promptConfirm, &confirmDelete); err != nil {
		return false, err
	}

	return confirmDelete, nil
}
