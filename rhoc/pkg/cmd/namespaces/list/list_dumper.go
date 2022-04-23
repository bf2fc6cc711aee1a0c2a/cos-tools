package list

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

type namespace struct {
	ID         string
	ClusterID  string
	Owner      string
	TenatKind  string
	TenatID    string
	State      string
	Expiration string
}
type namespaceWide struct {
	ID         string
	Name       string
	ClusterID  string
	Owner      string
	TenantKind string
	TenantID   string
	State      string
	Expiration string
}

func dumpAsTable(f *factory.Factory, items admin.ConnectorNamespaceList, wide bool) {
	r := make([]interface{}, 0, len(items.Items))

	for i := range items.Items {
		k := items.Items[i]

		if wide {
			r = append(r, namespaceWide{
				ClusterID:  k.ClusterId,
				ID:         k.Id,
				Name:       k.Name,
				Owner:      k.Owner,
				TenantKind: string(k.Tenant.Kind),
				TenantID:   k.Tenant.Id,
				State:      string(*&k.Status.State),
				Expiration: k.Expiration,
			})
		} else {
			r = append(r, namespace{
				ClusterID:  k.ClusterId,
				ID:         k.Id,
				Owner:      k.Owner,
				TenatKind:  string(k.Tenant.Kind),
				TenatID:    k.Tenant.Id,
				State:      string(*&k.Status.State),
				Expiration: k.Expiration,
			})
		}
	}

	t := dumper.NewTable(map[string]func(s string) tablewriter.Colors{
		"state":      statusCustomizer,
		"expiration": expirationCustomizer,
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

func expirationCustomizer(s string) tablewriter.Colors {
	if s == "" {
		return tablewriter.Colors{}
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return tablewriter.Colors{}
	}

	if time.Now().After(t) {
		return tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
	}

	return tablewriter.Colors{}
}
