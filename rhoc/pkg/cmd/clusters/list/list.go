package list

import (
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"net/http"
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName  = "list"
	CommandAlias = "ls"
)

type cluster struct {
	ID         string
	Owner      string
	CreatedAt  time.Time
	ModifiedAt time.Time
	State      string
}
type clusterWide struct {
	ID         string
	Name       string
	Owner      string
	CreatedAt  time.Time
	ModifiedAt time.Time
	State      string
}

type options struct {
	outputFormat string
	page         int
	limit        int
	all          bool
	orderBy      string
	search       string

	f *factory.Factory
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
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, cmdutil.ValidOutputs()...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, cmdutil.ValidOutputs()...)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddPage(cmd, &opts.page)
	cmdutil.AddLimit(cmd, &opts.limit)
	cmdutil.AddAllPages(cmd, &opts.all)
	cmdutil.AddOrderBy(cmd, &opts.orderBy)
	cmdutil.AddSearch(cmd, &opts.search)

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	items := admin.ConnectorClusterList{
		Kind:  "ConnectorClusterList",
		Items: make([]admin.ConnectorCluster, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.page; i == opts.page || opts.all; i++ {
		e := c.ConnectorClustersAdminApi.ListConnectorClusters(opts.f.Context)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.limit))

		if opts.orderBy != "" {
			e = e.OrderBy(opts.orderBy)
		}
		if opts.search != "" {
			e = e.Search(opts.search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			if httpRes != nil && httpRes.StatusCode == http.StatusInternalServerError {
				e, _ := service.ReadError(httpRes)
				if e.Reason != "" {
					err = fmt.Errorf("%s: [%w]", err.Error(), errors.New(e.Reason))
				}
			}
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			return err
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	if len(items.Items) == 0 && opts.outputFormat == "" {
		opts.f.Logger.Info("No result")
		return nil
	}

	switch opts.outputFormat {

	case dump.EmptyFormat:
		dumpAsTable(opts.f, items, false)
	case "wide":
		dumpAsTable(opts.f, items, true)
	default:
		return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, items)
	}

	return nil
}

func dumpAsTable(f *factory.Factory, items admin.ConnectorClusterList, wide bool) {
	r := make([]interface{}, 0, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		if wide {
			r = append(r, clusterWide{
				ID:         k.Id,
				Name:       k.Name,
				Owner:      k.Owner,
				CreatedAt:  k.CreatedAt,
				ModifiedAt: k.ModifiedAt,
				State:      string(k.Status.State),
			})
		} else {
			r = append(r, cluster{
				ID:         k.Id,
				Owner:      k.Owner,
				CreatedAt:  k.CreatedAt,
				ModifiedAt: k.ModifiedAt,
				State:      string(k.Status.State),
			})
		}
	}

	t := dumper.NewTable(map[string]func(s string) tablewriter.Colors{
		"state": statusCustomizer,
	})

	t.Dump(r, f.IOStreams.Out)
}

func statusCustomizer(s string) tablewriter.Colors {
	switch s {
	case "ready":
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
	case "disconnected":
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
	}

	return tablewriter.Colors{}
}
