package list

import (
	"errors"
	"fmt"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/request"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName  = "list"
	CommandAlias = "ls"
)

type options struct {
	request.ListDeploymentsOptions

	outputFormat string
	clusterID    string
	namespaceID  string

	f *factory.Factory
}

type deploymentDetail struct {
	admin.ConnectorDeploymentAdminView `json:",inline" yaml:",inline"`
	PlatformID                         string `json:"platform_id,omitempty" yaml:"platform_id,omitempty"`
}

func NewListCommand(f *factory.Factory) *cobra.Command {
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
			if opts.clusterID != "" && opts.namespaceID != "" {
				return errors.New("set either cluster-id or namespace-id, not both")
			}
			if opts.clusterID == "" && opts.namespaceID == "" {
				return errors.New("either cluster-id or namespace-id are required")
			}
			if opts.namespaceID != "" && (opts.DanglingDeployments || opts.ChannelUpdate) {
				return errors.New("dangling-deployments and channel-update are not supported with namespace-id at the moment")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddPage(cmd, &opts.Page)
	cmdutil.AddLimit(cmd, &opts.Limit)
	cmdutil.AddAllPages(cmd, &opts.AllPages)
	cmdutil.AddOrderBy(cmd, &opts.OrderBy)
	cmdutil.AddClusterID(cmd, &opts.clusterID)
	cmdutil.AddNamespaceID(cmd, &opts.namespaceID)
	cmdutil.AddChannelUpdate(cmd, &opts.ChannelUpdate)
	cmdutil.AddDanglingDeployments(cmd, &opts.DanglingDeployments)

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	var deployments admin.ConnectorDeploymentAdminViewList
	var cluster *admin.ConnectorClusterAdminView

	switch {
	case opts.clusterID != "":
		deployments, err = service.ListDeploymentsForCluster(c, opts.ListDeploymentsOptions, opts.clusterID)
		if err == nil {
			cluster, err = service.GetClusterByID(c, opts.clusterID)
		}
	case opts.namespaceID != "":
		deployments, err = service.ListDeploymentsForNamespace(c, opts.ListOptions, opts.namespaceID)
		if err == nil {
			cluster, err = service.GetClusterForNamespace(c, opts.namespaceID)
		}
	}

	if err != nil {
		return err
	}

	if len(deployments.Items) == 0 && opts.outputFormat == "" {
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "No result\n")
		return nil
	}

	items := make([]deploymentDetail, len(deployments.Items))
	for i := range deployments.Items {
		items[i] = deploymentDetail{
			ConnectorDeploymentAdminView: deployments.Items[i],
			PlatformID:                   cluster.Status.Platform.Id,
		}
	}

	switch opts.outputFormat {
	case dump.EmptyFormat:
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
		dumpAsTable(opts.f.IOStreams.Out, items, opts.ListDeploymentsOptions, false, dumper.TableStyleDefault)
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
	case cmdutil.OutputFormatWide:
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
		dumpAsTable(opts.f.IOStreams.Out, items, opts.ListDeploymentsOptions, true, dumper.TableStyleDefault)
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
	case cmdutil.OutputFormatCSV:
		dumpAsTable(opts.f.IOStreams.Out, items, opts.ListDeploymentsOptions, true, dumper.TableStyleCSV)
	default:
		return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, items)
	}

	return nil
}
