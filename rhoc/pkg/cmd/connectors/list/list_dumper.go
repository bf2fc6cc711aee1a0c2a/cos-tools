package list

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

type connector struct {
	ID              string
	NamespaceID     string
	Owner           string
	CreatedAt       time.Time
	ModifiedAt      time.Time
	ConnectorTypeId string
	Revision        int64
	DesiredState    string
	State           string
}

type connectorWide struct {
	ID              string
	Name            string
	NamespaceID     string
	Owner           string
	CreatedAt       time.Time
	ModifiedAt      time.Time
	ConnectorTypeId string
	Revision        int64
	DesiredState    string
	State           string
	Error           string
}

func dumpAsTable(f *factory.Factory, items admin.ConnectorAdminViewList, wide bool) {
	r := make([]interface{}, 0, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		if wide {
			r = append(r, connectorWide{
				NamespaceID:     k.NamespaceId,
				ID:              k.Id,
				Name:            k.Name,
				Owner:           k.Owner,
				CreatedAt:       k.CreatedAt,
				ModifiedAt:      k.ModifiedAt,
				ConnectorTypeId: k.ConnectorTypeId,
				Revision:        k.ResourceVersion,
				DesiredState:    string(k.DesiredState),
				State:           string(k.Status.State),
				Error:           k.Status.Error,
			})
		} else {
			r = append(r, connector{
				NamespaceID:     k.NamespaceId,
				ID:              k.Id,
				Owner:           k.Owner,
				CreatedAt:       k.CreatedAt,
				ModifiedAt:      k.ModifiedAt,
				ConnectorTypeId: k.ConnectorTypeId,
				Revision:        k.ResourceVersion,
				DesiredState:    string(k.DesiredState),
				State:           string(k.Status.State),
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
