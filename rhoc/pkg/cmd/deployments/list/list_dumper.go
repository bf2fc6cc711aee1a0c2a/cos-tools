package list

import (
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"k8s.io/apimachinery/pkg/util/duration"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorDeploymentAdminViewList, wide bool, csv bool) {
	t := dumper.Table[admin.ConnectorDeploymentAdminView]{
		Config: dumper.TableConfig{
			CSV:  csv,
			Wide: wide,
		},
		Columns: []dumper.Column[admin.ConnectorDeploymentAdminView]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "ConnectorID",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Spec.ConnectorId,
					}
				},
			},
			{
				Name: "NamespaceID",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Spec.NamespaceId,
					}
				},
			},
			{
				Name: "ClusterID",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Spec.ClusterId,
					}
				},
			},
			{
				Name: "ConnectorTypeId",
				Wide: true,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Spec.ConnectorTypeId,
					}
				},
			},
			{
				Name: "DesiredState",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: string(in.Spec.DesiredState),
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					r := dumper.Row{
						Value: string(in.Status.Phase),
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
				Name: "Version",
				Wide: true,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: strconv.FormatInt(in.Metadata.ResourceVersion, 10),
					}
				},
			},
			{
				Name: "DeploymentVersion",
				Wide: true,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					r := dumper.Row{
						Value: strconv.FormatInt(in.Status.ResourceVersion, 10),
					}

					switch {
					case in.Metadata.ResourceVersion > in.Status.ResourceVersion:
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor}

					case in.Metadata.ResourceVersion < in.Status.ResourceVersion:
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
					case in.Metadata.ResourceVersion == in.Status.ResourceVersion:
						r.Colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}
					}

					return r
				},
			},
			{
				Name: "Age",
				Wide: false,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					age := duration.HumanDuration(time.Since(in.Metadata.CreatedAt))
					if in.Metadata.CreatedAt.IsZero() {
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
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Metadata.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Metadata.UpdatedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	t.Dump(f.IOStreams.Out, items.Items)
}
