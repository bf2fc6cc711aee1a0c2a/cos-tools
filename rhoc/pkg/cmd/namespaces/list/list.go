package list

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
)

// row is the details of a Kafka instance needed to print to a table
type itemRow struct {
	ClusterID  string `json:"cluster_id" header:"ClusterID"`
	ID         string `json:"id" header:"ID"`
	Name       string `json:"name" header:"Name"`
	TenatKind  string `json:"tenant_kind" header:"TenantKind"`
	TenatID    string `json:"tenant_id" header:"TenantID"`
	Status     string `json:"status" header:"Status"`
	Expiration string `json:"expiration" header:"Expiration"`
}

type options struct {
	outputFormat string
	page         int
	limit        int
	all          bool
	clusterID    string
	orderBy      string
	search       string

	f *factory.Factory
}

func NewListCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list",
		Long:    "list",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, flagutil.ValidOutputFormats...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, flagutil.ValidOutputFormats...)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	flags := flagutil.NewFlagSet(cmd, f.Localizer)

	flags.AddOutput(&opts.outputFormat)
	flags.IntVarP(&opts.page, "page", "p", build.DefaultPageNumber, "page")
	flags.IntVarP(&opts.limit, "limit", "l", build.DefaultPageSize, "limit")
	flags.BoolVar(&opts.all, "all", false, "all")
	flags.StringVarP(&opts.clusterID, "cluster-id", "c", "", "cluster-id")
	flags.StringVar(&opts.orderBy, "order-by", "", "order-by")
	flags.StringVar(&opts.search, "search", "", "search")

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	items := admin.ConnectorNamespaceList{
		Kind:  "ConnectorNamespaceList",
		Items: make([]admin.ConnectorNamespace, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.page; i == opts.page || opts.all; i++ {

		var result *admin.ConnectorNamespaceList
		var err error
		var httpRes *http.Response

		if opts.clusterID != "" {
			e := c.ConnectorClustersAdminApi.GetClusterNamespaces(opts.f.Context, opts.clusterID)
			e = e.Page(strconv.Itoa(i))
			e = e.Size(strconv.Itoa(opts.limit))

			if opts.orderBy != "" {
				e = e.OrderBy(opts.orderBy)
			}
			if opts.search != "" {
				e = e.Search(opts.search)
			}

			result, httpRes, err = e.Execute()
		} else {
			e := c.ConnectorNamespacesAdminApi.GetConnectorNamespaces(opts.f.Context)
			e = e.Page(strconv.Itoa(i))
			e = e.Size(strconv.Itoa(opts.limit))

			if opts.orderBy != "" {
				e = e.OrderBy(opts.orderBy)
			}
			if opts.search != "" {
				e = e.Search(opts.search)
			}

			result, httpRes, err = e.Execute()
		}

		if httpRes != nil {
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
		rows := responseToRows(items)
		dump.Table(opts.f.IOStreams.Out, rows)
		opts.f.Logger.Info("")
	default:
		return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, items)
	}

	return nil
}

func responseToRows(items admin.ConnectorNamespaceList) []itemRow {
	rows := make([]itemRow, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		row := itemRow{
			ClusterID:  k.ClusterId,
			ID:         k.Id,
			Name:       k.Name,
			TenatKind:  string(k.Tenant.Kind),
			TenatID:    k.Tenant.Id,
			Status:     string(*&k.Status.State),
			Expiration: k.Expiration,
		}

		rows[i] = row
	}

	return rows
}
