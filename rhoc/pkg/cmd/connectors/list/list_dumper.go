package list

import (
	"io"
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/util/duration"
)

func dumpAsTable(out io.Writer, items []connectorDetail, wide bool, style dumper.TableStyle) {
	config := dumper.TableConfig[connectorDetail]{
		Style: style,
		Wide:  wide,
		Columns: []dumper.Column[connectorDetail]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "NamespaceId",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.NamespaceId,
					}
				},
			},
			{
				Name: "ClusterId",
				Wide: true,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.ClusterID,
					}
				},
			},
			{
				Name: "PlatformId",
				Wide: true,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.PlatformID,
					}
				},
			},
			{
				Name: "Name",
				Wide: true,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.Name,
					}
				},
			},
			{
				Name: "Owner",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.Owner,
					}
				},
			},
			{
				Name: "ConnectorTypeId",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.ConnectorTypeId,
					}
				},
			},
			{
				Name: "DesiredState",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: string(in.DesiredState),
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
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
				Name: "Reason",
				Wide: true,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.Status.Error,
					}
				},
			},
			{
				Name: "Version",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: strconv.FormatInt(in.ResourceVersion, 10),
					}
				},
			},
			{
				Name: "Age",
				Wide: false,
				Getter: func(in *connectorDetail) dumper.Row {
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
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *connectorDetail) dumper.Row {
					return dumper.Row{
						Value: in.ModifiedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	dumper.DumpTable(config, out, items)
}
