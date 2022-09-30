package get

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/ocm"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	ocmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/pkg/errors"
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

	// enables to retrieve info from ocm
	// ocm get /api/clusters_mgmt/v1/clusters --parameter=search="external_id like '...'"
	ocm bool

	f *factory.Factory
}

type ocmInfo struct {
	ID            string `json:"id,omitempty" yaml:"id,omitempty"`
	ClusterID     string `json:"cluster_id,omitempty" yaml:"cluster_id,omitempty"`
	Console       string `json:"cluster_console,omitempty" yaml:"cluster_console,omitempty"`
	ProductID     string `json:"product_id,omitempty" yaml:"product_id,omitempty"`
	CloudProvider string `json:"cloud_provider,omitempty" yaml:"cloud_provider,omitempty"`
}
type info struct {
	admin.ConnectorClusterAdminView `json:",inline" yaml:",inline"`
	Ocm                             *ocmInfo `json:"ocm,omitempty" yaml:"ocm,omitempty"`
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

	cmd.Flags().BoolVar(&opts.ocm, "ocm", opts.ocm, "ocm")

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	result, httpRes, err := c.Clusters().GetConnectorCluster(opts.f.Context, opts.id).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return response.Error(err, httpRes)
	}

	var cluster *ocmv1.Cluster

	if opts.ocm {
		cluster, err = ocm.Cluster(result.Status.Platform.Id)
		if err != nil {
			return errors.Wrap(err, "unable to retrieve cluster info")
		}
	}

	i := info{
		ConnectorClusterAdminView: *result,
	}

	if cluster != nil {
		i.Ocm = &ocmInfo{
			ID:        cluster.ExternalID(),
			ClusterID: cluster.ID(),
		}

		if cluster.Product() != nil {
			i.Ocm.ProductID = cluster.Product().ID()
		}
		if cluster.CloudProvider() != nil {
			i.Ocm.CloudProvider = cluster.CloudProvider().ID()
		}
		if cluster.Console() != nil {
			i.Ocm.Console = cluster.Console().URL()
		}
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, i)
}
