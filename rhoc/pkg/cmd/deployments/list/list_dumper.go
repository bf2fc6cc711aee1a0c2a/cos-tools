package list

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"time"
)

type deployment struct {
	ID              string
	ConnectorID     string
	NamespaceID     string
	ClusterID       string
	ResourceVersion int64
	CreatedAt       time.Time
	ModifiedAt      time.Time
	Status          string
}

type deploymentWide struct {
	ID              string
	ConnectorID     string
	NamespaceID     string
	ClusterID       string
	ResourceVersion int64
	OperatorID      string
	CreatedAt       time.Time
	ModifiedAt      time.Time
	Status          string
}

func dumpAsTable(f *factory.Factory, items admin.ConnectorDeploymentAdminViewList, wide bool) {
	r := make([]interface{}, 0, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		if wide {
			r = append(r, deploymentWide{
				ClusterID:       k.Spec.ClusterId,
				ConnectorID:     k.Spec.ConnectorId,
				NamespaceID:     k.Spec.NamespaceId,
				ID:              k.Id,
				ResourceVersion: k.Metadata.ResourceVersion,
				OperatorID:      k.Spec.OperatorId,
				CreatedAt:       k.Metadata.CreatedAt,
				ModifiedAt:      k.Metadata.UpdatedAt,
				Status:          string(*&k.Status.Phase),
			})
		} else {
			r = append(r, deployment{
				ClusterID:       k.Spec.ClusterId,
				ConnectorID:     k.Spec.ConnectorId,
				NamespaceID:     k.Spec.NamespaceId,
				ID:              k.Id,
				ResourceVersion: k.Metadata.ResourceVersion,
				CreatedAt:       k.Metadata.CreatedAt,
				ModifiedAt:      k.Metadata.UpdatedAt,
				Status:          string(*&k.Status.Phase),
			})
		}
	}

	t := dumper.NewTable(map[string]func(s string) tablewriter.Colors{
		"state": statusCustomizer,
	})

	t.Dump(r, f.IOStreams.Out)
}

func statusCustomizer(s string) tablewriter.Colors {
	switch s {
	case "ready":
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiGreenColor}
	case "failed":
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiRedColor}
	case "stopped":
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiYellowColor}
	}

	return tablewriter.Colors{}
}
