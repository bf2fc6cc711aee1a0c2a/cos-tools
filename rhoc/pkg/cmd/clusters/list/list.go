package list

import (
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

type itemRow struct {
	ID         string    `json:"id,omitempty" header:"ID"`
	Owner      string    `json:"owner,omitempty" header:"Owner"`
	CreatedAt  time.Time `json:"created_at,omitempty" header:"CreatedAt"`
	ModifiedAt time.Time `json:"modified_at,omitempty" header:"ModifiedAt"`
	Name       string    `json:"name,omitempty"  header:"Name"`
	Status     string    `json:"status" header:"Status"`
	Error      string    `json:"error" header:"Error"`
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

func responseToRows(items admin.ConnectorClusterList) []itemRow {
	rows := make([]itemRow, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		row := itemRow{
			ID:         k.Id,
			Owner:      k.Owner,
			Name:       k.Name,
			CreatedAt:  k.CreatedAt,
			ModifiedAt: k.ModifiedAt,
			Status:     string(k.Status.State),
			Error:      k.Status.Error,
		}

		rows[i] = row
	}

	return rows
}
