package list

import (
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

func dumpAsTable(f *factory.Factory, items admin.ConnectorAdminViewList, wide bool) {
	t := dumper.Table[admin.ConnectorAdminView]{}

	t.Field("ID", func(in *admin.ConnectorAdminView) string {
		return in.Id
	})

	if wide {
		t.Field("Name", func(in *admin.ConnectorAdminView) string {
			return in.Name
		})
	}

	t.Field("NamespaceID", func(in *admin.ConnectorAdminView) string {
		return in.NamespaceId
	})

	t.Field("Owner", func(in *admin.ConnectorAdminView) string {
		return in.Owner
	})

	t.Field("ConnectorTypeId", func(in *admin.ConnectorAdminView) string {
		return in.ConnectorTypeId
	})

	if wide {
		t.Field("CreatedAt", func(in *admin.ConnectorAdminView) string {
			return in.CreatedAt.Format(time.RFC3339)
		})

		t.Field("ModifiedAt", func(in *admin.ConnectorAdminView) string {
			return in.ModifiedAt.Format(time.RFC3339)
		})

		t.Field("ResourceVersion", func(in *admin.ConnectorAdminView) string {
			return strconv.FormatInt(in.ResourceVersion, 10)
		})
	}

	t.Field("DesiredState", func(in *admin.ConnectorAdminView) string {
		return string(in.DesiredState)
	})

	t.Rich("State", func(in *admin.ConnectorAdminView) (string, tablewriter.Colors) {
		s := string(in.Status.State)
		c := tablewriter.Colors{}

		switch s {
		case "ready":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
		case "failed":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiRedColor}
		case "stopped":
			c = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiYellowColor}
		}

		return s, c
	})

	t.Dump(items.Items, f.IOStreams.Out)
}
