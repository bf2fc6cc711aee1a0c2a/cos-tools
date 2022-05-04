package list

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"k8s.io/apimachinery/pkg/util/duration"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorNamespaceList, wide bool, csv bool) {
	t := dumper.Table[admin.ConnectorNamespace]{
		Config: dumper.TableConfig{
			CSV:  csv,
			Wide: wide,
		},
		Columns: []dumper.Column[admin.ConnectorNamespace]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "ClusterID",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.ClusterId,
					}
				},
			},
			{
				Name: "Name",
				Wide: true,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.Name,
					}
				},
			},
			{
				Name: "Owner",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.Owner,
					}
				},
			},
			{
				Name: "TenantKind",
				Wide: true,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: string(in.Tenant.Kind),
					}
				},
			},

			{
				Name: "TenantID",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.Tenant.Id,
					}
				},
			},
			{
				Name: "Version",
				Wide: true,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: strconv.FormatInt(in.ResourceVersion, 10),
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					r := dumper.Row{
						Value: string(in.Status.State),
					}

					switch r.Value {
					case "ready":
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
					case "disconnected":
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
					}

					if wide && in.Expiration != "" {
						t, err := time.Parse(time.RFC3339, in.Expiration)
						if err == nil && time.Now().After(t) {
							r.Value = r.Value + " (*)"
						}
					}

					return r
				},
			},
			{
				Name: "Connectors",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: fmt.Sprint(in.Status.ConnectorsDeployed),
					}
				},
			},
			{
				Name: "Expiration",
				Wide: true,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {

					r := dumper.Row{
						Value: in.Expiration,
					}

					if r.Value != "" {
						t, err := time.Parse(time.RFC3339, r.Value)
						if err == nil && time.Now().After(t) {
							r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
						}
					}

					return r
				},
			},
			{
				Name: "Age",
				Wide: false,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
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
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *admin.ConnectorNamespace) dumper.Row {
					return dumper.Row{
						Value: in.ModifiedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	t.Dump(f.IOStreams.Out, items.Items)
}
