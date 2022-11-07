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
	request.ListOptions

	outputFormat string
	clusterID    string
	namespaceID  string

	f *factory.Factory
}

func NewListCommand(f *factory.Factory) *cobra.Command {

	opts := options{
		f: f,
	}

	sc, err := f.ServiceContext.Load()
	if err == nil && sc != nil {
		opts.namespaceID = sc.Contexts[sc.CurrentContext].NamespaceID
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
	cmdutil.AddNamespaceID(cmd, &opts.namespaceID)

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	var items admin.ConnectorAdminViewList

	switch {
	case opts.clusterID != "":
		items, err = service.ListConnectorsForCluster(c, opts.ListOptions, opts.clusterID)
	case opts.namespaceID != "":
		items, err = service.ListConnectorsForNamespace(c, opts.ListOptions, opts.namespaceID)
	}

	if err != nil {
		return err
	}

	if len(items.Items) == 0 && opts.outputFormat == "" {
		_, _ = fmt.Fprint(opts.f.IOStreams.Out, "No result\n")
		return nil
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
