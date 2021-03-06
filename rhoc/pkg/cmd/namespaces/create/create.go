package create

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/connector/connectorcmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName = "create"
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
		Use:  CommandName,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.ValidateOutputs(cmd); err != nil {
				return err
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

	e := c.Namespaces().CreateConnectorNamespace(opts.f.Context)
	e = e.ConnectorNamespaceWithTenantRequest(r)

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
