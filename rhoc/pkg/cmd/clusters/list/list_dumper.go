package list

import (
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"k8s.io/apimachinery/pkg/util/duration"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorClusterList, wide bool, csv bool) {
	t := dumper.Table[admin.ConnectorCluster]{
		Config: dumper.TableConfig{
			CSV:  csv,
			Wide: wide,
		},
	}

	t.Column("ID", false, func(in *admin.ConnectorCluster) dumper.Row {
		return dumper.Row{
			Value: in.Id,
		}
	})

	t.Column("Name", true, func(in *admin.ConnectorCluster) dumper.Row {
		return dumper.Row{
			Value: in.Name,
		}
	})

	t.Column("Owner", false, func(in *admin.ConnectorCluster) dumper.Row {
		return dumper.Row{
			Value: in.Owner,
		}
	})

	t.Column("State", false, func(in *admin.ConnectorCluster) dumper.Row {
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
	})

	t.Column("Age", false, func(in *admin.ConnectorCluster) dumper.Row {
		age := duration.HumanDuration(time.Since(in.CreatedAt))
		if in.CreatedAt.IsZero() {
			age = ""
		}

		return dumper.Row{
			Value: age,
		}
	})

	t.Column("CreatedAt", true, func(in *admin.ConnectorCluster) dumper.Row {
		return dumper.Row{
			Value: in.CreatedAt.Format(time.RFC3339),
		}
	})

	t.Column("ModifiedAt", true, func(in *admin.ConnectorCluster) dumper.Row {
		return dumper.Row{
			Value: in.ModifiedAt.Format(time.RFC3339),
		}
	})

	t.Dump(f.IOStreams.Out, items.Items)
}
