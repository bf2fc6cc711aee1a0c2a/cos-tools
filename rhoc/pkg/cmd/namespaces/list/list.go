package list

import (
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
	request.ListOptions

	outputFormat string
	clusterID    string

	f *factory.Factory
}

type namespaceDetail struct {
	admin.ConnectorNamespace `json:",inline" yaml:",inline"`
	PlatformID               string `json:"platform_id,omitempty" yaml:"platform_id,omitempty"`
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
	cmdutil.AddSearch(cmd, &opts.Search)
	cmdutil.AddClusterID(cmd, &opts.clusterID)

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	var namespaces admin.ConnectorNamespaceList

	switch {
	case opts.clusterID != "":
		namespaces, err = service.ListNamespacesForCluster(c, opts.ListOptions, opts.clusterID)
	default:
		namespaces, err = service.ListNamespaces(c, opts.ListOptions)
	}

	if err != nil {
		return err
	}

	if len(namespaces.Items) == 0 && opts.outputFormat == "" {
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "No result\n")
		return nil
	}

	clusters := make(map[string]*admin.ConnectorClusterAdminView)

	items := make([]namespaceDetail, len(namespaces.Items))
	for i := range namespaces.Items {
		cluster := clusters[namespaces.Items[i].ClusterId]
		if cluster == nil {
			cluster, err = service.GetClusterByID(c, namespaces.Items[i].ClusterId)
			if err != nil {
				return err
			}

			clusters[namespaces.Items[i].ClusterId] = cluster
		}

		items[i] = namespaceDetail{
			ConnectorNamespace: namespaces.Items[i],
			PlatformID:         cluster.Status.Platform.Id,
		}
	}

	switch opts.outputFormat {
	case dump.EmptyFormat:
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
		dumpAsTable(opts.f.IOStreams.Out, items, false, dumper.TableStyleDefault)
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
	case cmdutil.OutputFormatWide:
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
		dumpAsTable(opts.f.IOStreams.Out, items, true, dumper.TableStyleDefault)
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "\n")
	case cmdutil.OutputFormatCSV:
		dumpAsTable(opts.f.IOStreams.Out, items, true, dumper.TableStyleCSV)
	default:
		return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, items)
	}

	return nil
}
