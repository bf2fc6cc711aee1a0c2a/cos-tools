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

func dumpAsTable(f *factory.Factory, items admin.ConnectorDeploymentAdminViewList, wide bool) {
	t := dumper.Table[admin.ConnectorDeploymentAdminView]{}

	t.Field("ID", func(in *admin.ConnectorDeploymentAdminView) string {
		return in.Id
	})

	t.Field("NamespaceID", func(in *admin.ConnectorDeploymentAdminView) string {
		return in.Spec.NamespaceId
	})

	t.Field("ClusterId", func(in *admin.ConnectorDeploymentAdminView) string {
		return in.Spec.ClusterId
	})

	t.Field("ConnectorTypeId", func(in *admin.ConnectorDeploymentAdminView) string {
		return in.Spec.ConnectorTypeId
	})

	if wide {

		t.Field("Version", func(in *admin.ConnectorDeploymentAdminView) string {
			return strconv.FormatInt(in.Metadata.ResourceVersion, 10)
		})

		t.Field("ConnectorVersion", func(in *admin.ConnectorDeploymentAdminView) string {
			return strconv.FormatInt(in.Spec.ConnectorResourceVersion, 10)
		})
	}

	t.Field("DesiredState", func(in *admin.ConnectorDeploymentAdminView) string {
		return string(in.Spec.DesiredState)
	})

	t.Rich("State", func(in *admin.ConnectorDeploymentAdminView) (string, tablewriter.Colors) {
		s := string(in.Status.Phase)
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

	if wide {
		t.Field("CreatedAt", func(in *admin.ConnectorDeploymentAdminView) string {
			return in.Metadata.CreatedAt.Format(time.RFC3339)
		})
		t.Field("UpdatedAt", func(in *admin.ConnectorDeploymentAdminView) string {
			return in.Metadata.UpdatedAt.Format(time.RFC3339)
		})
	} else {
		t.Field("Age", func(in *admin.ConnectorDeploymentAdminView) string {
			age := duration.HumanDuration(time.Since(in.Metadata.CreatedAt))
			if in.Metadata.CreatedAt.IsZero() {
				age = ""
			}

			return age
		})
	}

	t.Dump(f.IOStreams.Out, items.Items)
}
