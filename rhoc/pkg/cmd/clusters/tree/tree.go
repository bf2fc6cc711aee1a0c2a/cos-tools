package tree

import (
	"fmt"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/request"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/resource"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName = "tree"

	firstElemPrefix = `├─`
	lastElemPrefix  = `└─`
	indent          = "  "
	pipe            = `│ `
)

type options struct {
	request.ListOptions
	id string

	f *factory.Factory
}

func NewTreeCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:  CommandName,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddPage(cmd, &opts.Page)
	cmdutil.AddLimit(cmd, &opts.Limit)
	cmdutil.AddAllPages(cmd, &opts.AllPages)
	cmdutil.AddOrderBy(cmd, &opts.OrderBy)
	cmdutil.AddSearch(cmd, &opts.Search)
	cmdutil.AddID(cmd, &opts.id)

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(opts.f.IOStreams.Out)
	table.SetHeader([]string{"ID", "OWNER", "AGE", "STATUS", "REASON"})
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetColumnSeparator(tablewriter.SPACE)
	table.SetCenterSeparator(tablewriter.SPACE)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	if err := renderCluster(c, opts, table); err != nil {
		return err
	}

	namespaces, err := service.ListNamespacesForCluster(c, opts.ListOptions, opts.id)
	if err != nil {
		return err
	}

	for nsIndex := range namespaces.Items {
		if err := renderNamespace(namespaces, nsIndex, table); err != nil {
			return err
		}

		connectors, err := service.ListConnectorsForNamespace(c, opts.ListOptions, namespaces.Items[nsIndex].Id)
		if err != nil {
			return err
		}

		for cntIndex := range connectors.Items {
			if err := renderConnector(connectors, cntIndex, table); err != nil {
				return err
			}
		}
	}

	table.Render()

	return nil
}

func renderCluster(c *service.AdminAPI, opts *options, table *tablewriter.Table) error {
	cluster, httpRes, err := c.Clusters().GetConnectorCluster(c.Context(), opts.id).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return err
	}

	data := []string{
		"cluster/" + cluster.Id,
		cluster.Owner,
		resource.Age(cluster.CreatedAt),
		string(cluster.Status.State),
		cluster.Status.Error}

	style := []tablewriter.Colors{{}, {}, {}, {}, {}}

	switch string(cluster.Status.State) {
	case "ready":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
	case "disconnected":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
	}

	table.Rich(data, style)

	return nil
}

func renderNamespace(namespaces admin.ConnectorNamespaceList, index int, table *tablewriter.Table) error {
	var data []string

	ns := namespaces.Items[index]
	style := []tablewriter.Colors{{}, {}, {}, {}, {}}

	if index == len(namespaces.Items)-1 {
		data = []string{
			fmt.Sprintf("%s%s", lastElemPrefix, "namespace/"+ns.Id),
			ns.Owner,
			resource.Age(ns.CreatedAt),
			string(ns.Status.State),
			ns.Status.Error,
		}
	} else {
		data = []string{
			fmt.Sprintf("%s%s", firstElemPrefix, "namespace/"+ns.Id),
			ns.Owner,
			resource.Age(ns.CreatedAt),
			string(ns.Status.State),
			ns.Status.Error,
		}
	}

	switch ns.Tenant.Kind {
	case admin.CONNECTORNAMESPACETENANTKIND_USER:
		style[1] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor}
	case admin.CONNECTORNAMESPACETENANTKIND_ORGANISATION:
		style[1] = tablewriter.Colors{}
	}

	switch string(ns.Status.State) {
	case "ready":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
	case "disconnected":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
	}

	table.Rich(data, style)

	return nil
}

func renderConnector(connectors admin.ConnectorAdminViewList, index int, table *tablewriter.Table) error {
	var data []string

	ct := connectors.Items[index]
	style := []tablewriter.Colors{{}, {}, {}, {}, {}}

	if index == len(connectors.Items)-1 {
		data = []string{
			fmt.Sprintf("%s%s%s%s", pipe, indent, lastElemPrefix, "connector/"+ct.Id),
			ct.Owner,
			resource.Age(ct.CreatedAt),
			string(ct.Status.State),
			ct.Status.Error,
		}
	} else {
		data = []string{
			fmt.Sprintf("%s%s%s%s", pipe, indent, firstElemPrefix, "connector/"+ct.Id),
			ct.Owner,
			resource.Age(ct.CreatedAt),
			string(ct.Status.State),
			ct.Status.Error,
		}
	}

	switch string(ct.Status.State) {
	case "ready":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
	case "failed":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiRedColor}
	case "stopped":
		style[3] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiYellowColor}
	}

	table.Rich(data, style)

	return nil
}
