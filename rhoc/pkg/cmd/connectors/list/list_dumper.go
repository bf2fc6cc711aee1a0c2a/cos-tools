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
		Columns: []dumper.Column[admin.ConnectorAdminView]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "NamespaceId",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.NamespaceId,
					}
				},
			},
			{
				Name: "Name",
				Wide: true,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Name,
					}
				},
			},
			{
				Name: "Owner",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Owner,
					}
				},
			},
			{
				Name: "ConnectorTypeId",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.ConnectorTypeId,
					}
				},
			},
			{
				Name: "Version",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: strconv.FormatInt(in.ResourceVersion, 10),
					}
				},
			},
			{
				Name: "DesiredState",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: string(in.DesiredState),
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
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
				},
			},
			{
				Name: "Age",
				Wide: false,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
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
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *admin.ConnectorAdminView) dumper.Row {
					return dumper.Row{
						Value: in.ModifiedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	t.Dump(f.IOStreams.Out, items.Items)
}
