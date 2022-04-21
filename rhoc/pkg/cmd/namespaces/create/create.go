package create

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/connector/connectorcmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
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

	flags := flagutil.NewFlagSet(cmd, f.Localizer)
	flags.StringVar(&opts.name, "name", "", "name")
	flags.StringVar(&opts.clusterID, "cluster-id", "", "cluster-id")
	flags.StringVar(&opts.tenantKind, "tenant-kind", "", "tenant-kind")
	flags.StringVar(&opts.tenantID, "tenant-id", "", "tenant-id")
	flags.AddOutput(&opts.outputFormat)

	cmd.MarkFlagRequired("cluster-id")
	cmd.MarkFlagRequired("tenant-kind")
	cmd.MarkFlagRequired("tenant-id")

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
		return err
	}
	if httpRes != nil && httpRes.StatusCode == 204 {
		return nil
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, result)
}
