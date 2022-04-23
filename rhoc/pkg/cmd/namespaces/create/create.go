package create

import (
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/connector/connectorcmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"net/http"
)

type options struct {
	f *factory.Factory

	name         string
	outputFormat string
	clusterID    string
	tenantKind   string
	tenantID     string
}

func NewCreateCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create",
		Long:  "create",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, cmdutil.ValidOutputs()...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, cmdutil.ValidOutputs()...)
			}

			validator := connectorcmdutil.Validator{
				Localizer: f.Localizer,
			}

			if err := validator.ValidateNamespace(opts.name); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddName(cmd, &opts.name)
	cmdutil.AddClusterID(cmd, &opts.clusterID).Required()
	cmdutil.AddTenantKind(cmd, &opts.tenantKind).Required()
	cmdutil.AddTenantID(cmd, &opts.tenantID).Required()

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	r := admin.ConnectorNamespaceWithTenantRequest{
		Name:      opts.name,
		ClusterId: opts.clusterID,
		Tenant: admin.ConnectorNamespaceTenant{
			Id:   opts.tenantID,
			Kind: admin.ConnectorNamespaceTenantKind(opts.tenantKind),
		},
	}

	e := c.ConnectorNamespacesAdminApi.CreateConnectorNamespace(opts.f.Context)
	e = e.ConnectorNamespaceWithTenantRequest(r)

	result, httpRes, err := e.Execute()
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

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, result)
}
