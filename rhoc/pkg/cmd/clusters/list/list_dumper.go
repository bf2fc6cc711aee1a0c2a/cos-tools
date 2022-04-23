package list

import (
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorClusterList, wide bool) {
	t := dumper.Table[admin.ConnectorCluster]{}

	t.Field("ID", func(in *admin.ConnectorCluster) string {
		return in.Id
	})

	if wide {
		t.Field("Name", func(in *admin.ConnectorCluster) string {
			return in.Name
		})
	}

	t.Field("Owner", func(in *admin.ConnectorCluster) string {
		return in.Owner
	})

	if wide {
		t.Field("CreatedAt", func(in *admin.ConnectorCluster) string {
			return in.CreatedAt.Format(time.RFC3339)
		})

		t.Field("ModifiedAt", func(in *admin.ConnectorCluster) string {
			return in.ModifiedAt.Format(time.RFC3339)
		})
	}

	t.Rich("State", func(in *admin.ConnectorCluster) (string, tablewriter.Colors) {
		s := string(in.Status.State)
		c := tablewriter.Colors{}

		switch s {
		case "ready":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
		case "disconnected":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
		}

		return s, c
	})

	t.Dump(items.Items, f.IOStreams.Out)
}
