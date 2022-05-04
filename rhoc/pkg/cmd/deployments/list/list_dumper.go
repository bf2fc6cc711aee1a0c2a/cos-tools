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
	}

	t.Column("ID", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Id,
		}
	})

	t.Column("NamespaceID", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Spec.NamespaceId,
		}
	})

	t.Column("ClusterId", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Spec.ClusterId,
		}
	})

	t.Column("ConnectorTypeId", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Spec.ConnectorTypeId,
		}
	})

	t.Column("Version", true, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: strconv.FormatInt(in.Metadata.ResourceVersion, 10),
		}
	})

	t.Column("ConnectorVersion", true, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: strconv.FormatInt(in.Spec.ConnectorResourceVersion, 10),
		}
	})

	t.Column("DesiredState", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: string(in.Spec.DesiredState),
		}
	})

	t.Column("State", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
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
	})

	t.Column("Age", false, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		age := duration.HumanDuration(time.Since(in.Metadata.CreatedAt))
		if in.Metadata.CreatedAt.IsZero() {
			age = ""
		}

		return dumper.Row{
			Value: age,
		}
	})

	t.Column("CreatedAt", true, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Metadata.CreatedAt.Format(time.RFC3339),
		}
	})

	t.Column("ModifiedAt", true, func(in *admin.ConnectorDeploymentAdminView) dumper.Row {
		return dumper.Row{
			Value: in.Metadata.UpdatedAt.Format(time.RFC3339),
		}
	})

	t.Dump(f.IOStreams.Out, items.Items)
}
