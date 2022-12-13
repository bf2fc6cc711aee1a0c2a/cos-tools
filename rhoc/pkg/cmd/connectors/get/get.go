package get

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName = "get"
)

type options struct {
	id           string
	outputFormat string

	f *factory.Factory
}

type connectorDetail struct {
	admin.ConnectorAdminView `json:",inline" yaml:",inline"`
	ClusterID                string `json:"cluster_id,omitempty" yaml:"cluster_id,omitempty"`
	PlatformID               string `json:"platform_id,omitempty" yaml:"platform_id,omitempty"`
}

func NewGetCommand(f *factory.Factory) *cobra.Command {
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
			return run(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddID(cmd, &opts.id).Required()

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	connector, err := service.GetConnectorByID(c, opts.id)
	if err != nil {
		return err
	}

	cluster, err := service.GetClusterForNamespace(c, connector.NamespaceId)
	if err != nil {
		return err
	}

	detail := connectorDetail{
		ConnectorAdminView: *connector,
		ClusterID:          cluster.Id,
		PlatformID:         cluster.Status.Platform.Id,
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, detail)
}
