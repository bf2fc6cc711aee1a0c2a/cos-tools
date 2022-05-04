package list

import (
	"fmt"
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
	}

	t.Column("ID", false, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.Id,
		}
	})

	t.Column("Name", true, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.Name,
		}
	})

	t.Column("ClusterID", false, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.ClusterId,
		}
	})

	t.Column("Owner", false, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.Owner,
		}
	})

	t.Column("TenantKind", true, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: string(in.Tenant.Kind),
		}
	})

	t.Column("TenantID", false, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.Tenant.Id,
		}
	})

	t.Column("State", false, func(in *admin.ConnectorNamespace) dumper.Row {
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
	})

	t.Column("Connectors", false, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: fmt.Sprint(in.Status.ConnectorsDeployed),
		}
	})

	t.Column("Expiration", true, func(in *admin.ConnectorNamespace) dumper.Row {
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
	})

	t.Column("Age", false, func(in *admin.ConnectorNamespace) dumper.Row {
		age := duration.HumanDuration(time.Since(in.CreatedAt))
		if in.CreatedAt.IsZero() {
			age = ""
		}

		return dumper.Row{
			Value: age,
		}
	})

	t.Column("CreatedAt", true, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.CreatedAt.Format(time.RFC3339),
		}
	})

	t.Column("ModifiedAt", true, func(in *admin.ConnectorNamespace) dumper.Row {
		return dumper.Row{
			Value: in.ModifiedAt.Format(time.RFC3339),
		}
	})

	t.Dump(f.IOStreams.Out, items.Items)
}
