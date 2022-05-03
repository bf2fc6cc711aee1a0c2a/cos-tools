package tree

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/duration"
)

const (
	CommandName = "tree"

	firstElemPrefix = `├─`
	lastElemPrefix  = `└─`
	indent          = "  "
	pipe            = `│ `
)

type options struct {
	page    int
	limit   int
	all     bool
	id      string
	orderBy string
	search  string

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

	cmdutil.AddPage(cmd, &opts.page)
	cmdutil.AddLimit(cmd, &opts.limit)
	cmdutil.AddAllPages(cmd, &opts.all)
	cmdutil.AddOrderBy(cmd, &opts.orderBy)
	cmdutil.AddSearch(cmd, &opts.search)
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

	namespaces, err := listNamespaces(c, opts, opts.id)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(opts.f.IOStreams.Out)
	table.SetHeader([]string{"ID", "STATUS", "REASON", "AGE"})
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(false)
	table.SetColumnSeparator(tablewriter.SPACE)
	table.SetCenterSeparator(tablewriter.SPACE)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{opts.id, "", ""})

	for i, ns := range namespaces {
		age := duration.HumanDuration(time.Since(ns.CreatedAt))
		if ns.CreatedAt.IsZero() {
			age = ""
		}

		if i == len(namespaces)-1 {
			table.Append([]string{
				fmt.Sprintf("%s%s (%d)", lastElemPrefix, ns.Id, ns.Status.ConnectorsDeployed),
				string(ns.Status.State),
				ns.Status.Error,
				age,
			})
		} else {
			table.Append([]string{
				fmt.Sprintf("%s%s (%d)", firstElemPrefix, ns.Id, ns.Status.ConnectorsDeployed),
				string(ns.Status.State),
				ns.Status.Error,
				age,
			})
		}

		connectors, err := listConnectors(c, opts, ns.Id)
		if err != nil {
			return err
		}

		for i, ct := range connectors {
			data := []string{}
			style := []tablewriter.Colors{{}, {}, {}}

			age := duration.HumanDuration(time.Since(ct.CreatedAt))
			if ct.CreatedAt.IsZero() {
				age = ""
			}

			if i == len(connectors)-1 {
				data = []string{
					fmt.Sprintf("%s%s%s%s", pipe, indent, lastElemPrefix, ct.Id),
					string(ct.Status.State),
					ns.Status.Error,
					age,
				}
			} else {
				data = []string{
					fmt.Sprintf("%s%s%s%s", pipe, indent, firstElemPrefix, ct.Id),
					string(ct.Status.State),
					ns.Status.Error,
					age,
				}
			}

			switch string(ct.Status.State) {
			case "ready":
				style[1] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
			case "failed":
				style[1] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiRedColor}
			case "stopped":
				style[1] = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiYellowColor}
			}

			table.Rich(data, style)
		}
	}

	table.Render()

	return nil
}

func listNamespaces(c service.AdminAPI, opts *options, clusterId string) ([]admin.ConnectorNamespace, error) {
	items := make([]admin.ConnectorNamespace, 0)

	for i := opts.page; i == opts.page || opts.all; i++ {
		var result *admin.ConnectorNamespaceList
		var err error
		var httpRes *http.Response

		e := c.Clusters().GetClusterNamespaces(opts.f.Context, clusterId)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.limit))

		if opts.orderBy != "" {
			e = e.OrderBy(opts.orderBy)
		}
		if opts.search != "" {
			e = e.Search(opts.search)
		}

		result, httpRes, err = e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			return []admin.ConnectorNamespace{}, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items = append(items, result.Items...)
	}

	return items, nil
}

func listConnectors(c service.AdminAPI, opts *options, namespaceId string) ([]admin.ConnectorAdminView, error) {
	items := make([]admin.ConnectorAdminView, 0)

	for i := opts.page; i == opts.page || opts.all; i++ {
		var result *admin.ConnectorAdminViewList
		var err error
		var httpRes *http.Response

		e := c.Clusters().GetNamespaceConnectors(opts.f.Context, namespaceId)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.limit))

		if opts.orderBy != "" {
			e = e.OrderBy(opts.orderBy)
		}
		if opts.search != "" {
			e = e.Search(opts.search)
		}

		result, httpRes, err = e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			return []admin.ConnectorAdminView{}, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items = append(items, result.Items...)
	}

	return items, nil
}
