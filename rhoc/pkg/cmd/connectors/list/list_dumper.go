package list

import (
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"k8s.io/apimachinery/pkg/util/duration"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorAdminViewList, wide bool, csv bool) {
	t := dumper.Table[admin.ConnectorAdminView]{
		Config: dumper.TableConfig{
			CSV:  csv,
			Wide: wide,
		},
	}

	t.Column("ID", false, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Id,
		}
	})

	t.Column("Name", true, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Name,
		}
	})

	t.Column("NamespaceId", false, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.NamespaceId,
		}
	})

	t.Column("Owner", false, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Owner,
		}
	})

	t.Column("ConnectorTypeId", false, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.ConnectorTypeId,
		}
	})

	t.Column("Version", true, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: strconv.FormatInt(in.ResourceVersion, 10),
		}
	})

	t.Column("DesiredState", false, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: string(in.DesiredState),
		}
	})

	t.Column("State", false, func(in *admin.ConnectorAdminView) dumper.Row {
		r := dumper.Row{
			Value: string(in.Status.State),
		}

		switch r.Value {
		case "ready":
			r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
		case "failed":
			r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiRedColor}
		case "stopped":
			r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiYellowColor}
		}

		return r
	})

	t.Column("Age", false, func(in *admin.ConnectorAdminView) dumper.Row {
		age := duration.HumanDuration(time.Since(in.CreatedAt))
		if in.CreatedAt.IsZero() {
			age = ""
		}

		return dumper.Row{
			Value: age,
		}
	})

	t.Column("CreatedAt", true, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.CreatedAt.Format(time.RFC3339),
		}
	})

	t.Column("ModifiedAt", true, func(in *admin.ConnectorAdminView) dumper.Row {
		return dumper.Row{
			Value: in.ModifiedAt.Format(time.RFC3339),
		}
	})

	t.Dump(f.IOStreams.Out, items.Items)
}
