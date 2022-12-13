package list

import (
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"io"
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/util/duration"
)

func dumpAsTable(out io.Writer, items []namespaceDetail, wide bool, style dumper.TableStyle) {
	config := dumper.TableConfig[namespaceDetail]{
		Style: style,
		Wide:  wide,
		Columns: []dumper.Column[namespaceDetail]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "ClusterID",
				Wide: false,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.ClusterId,
					}
				},
			},
			{
				Name: "PlatformID",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.PlatformID,
					}
				},
			},
			{
				Name: "Name",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.Name,
					}
				},
			},
			{
				Name: "Owner",
				Wide: false,
				Getter: func(in *namespaceDetail) dumper.Row {
					r := dumper.Row{
						Value: in.Owner,
					}

					switch in.Tenant.Kind {
					case admin.CONNECTORNAMESPACETENANTKIND_USER:
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor}
					}

					return r
				},
			},
			{
				Name: "TenantKind",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: string(in.Tenant.Kind),
					}
				},
			},

			{
				Name: "TenantID",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.Tenant.Id,
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *namespaceDetail) dumper.Row {
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
				Name: "Version",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: strconv.FormatInt(in.ResourceVersion, 10),
					}
				},
			},
			{
				Name: "Connectors",
				Wide: false,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: fmt.Sprint(in.Status.ConnectorsDeployed),
					}
				},
			},
			{
				Name: "Expiration",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {

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
				Getter: func(in *namespaceDetail) dumper.Row {
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
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *namespaceDetail) dumper.Row {
					return dumper.Row{
						Value: in.ModifiedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	dumper.DumpTable(config, out, items)
}
