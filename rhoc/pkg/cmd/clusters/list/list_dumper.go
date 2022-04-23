package list

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"time"
)

type cluster struct {
	ID         string
	Owner      string
	CreatedAt  time.Time
	ModifiedAt time.Time
	State      string
}
type clusterWide struct {
	ID         string
	Name       string
	Owner      string
	CreatedAt  time.Time
	ModifiedAt time.Time
	State      string
}

func dumpAsTable(f *factory.Factory, items admin.ConnectorClusterList, wide bool) {
	r := make([]interface{}, 0, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		if wide {
			r = append(r, clusterWide{
				ID:         k.Id,
				Name:       k.Name,
				Owner:      k.Owner,
				CreatedAt:  k.CreatedAt,
				ModifiedAt: k.ModifiedAt,
				State:      string(k.Status.State),
			})
		} else {
			r = append(r, cluster{
				ID:         k.Id,
				Owner:      k.Owner,
				CreatedAt:  k.CreatedAt,
				ModifiedAt: k.ModifiedAt,
				State:      string(k.Status.State),
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
	case "disconnected":
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}
	}

	return tablewriter.Colors{}
}
