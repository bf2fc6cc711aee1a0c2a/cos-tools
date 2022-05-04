package list

import (
	"io"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/util/duration"
)

func dumpAsTable(out io.Writer, items admin.ConnectorClusterList, wide bool, csv bool) {
	t := dumper.Table[admin.ConnectorCluster]{
		Config: dumper.TableConfig{
			CSV:  csv,
			Wide: wide,
		},
		Columns: []dumper.Column[admin.ConnectorCluster]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "Name",
				Wide: true,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					return dumper.Row{
						Value: in.Name,
					}
				},
			},
			{
				Name: "Owner",
				Wide: false,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					return dumper.Row{
						Value: in.Owner,
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					r := dumper.Row{
						Value: string(in.Status.State),
					}

					switch r.Value {
					case "ready":
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
					case "disconnected":
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
					}

					return r
				},
			},
			{
				Name: "Age",
				Wide: false,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					age := duration.HumanDuration(time.Since(in.CreatedAt))
					if in.CreatedAt.IsZero() {
						age = ""
					}

					return dumper.Row{
						Value: age,
					}
				},
			},
			{
				Name: "CreatedAt",
				Wide: true,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					return dumper.Row{
						Value: in.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *admin.ConnectorCluster) dumper.Row {
					return dumper.Row{
						Value: in.ModifiedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	t.Dump(out, items.Items)
}
