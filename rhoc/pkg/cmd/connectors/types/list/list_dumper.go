package list

import (
	"io"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
)

func dumpAsTable(out io.Writer, items admin.ConnectorTypeAdminViewList, wide bool, style dumper.TableStyle) {
	config := dumper.TableConfig[admin.ConnectorTypeAdminView]{
		Style: style,
		Wide:  wide,
		Columns: []dumper.Column[admin.ConnectorTypeAdminView]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *admin.ConnectorTypeAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "Name",
				Wide: true,
				Getter: func(in *admin.ConnectorTypeAdminView) dumper.Row {
					return dumper.Row{
						Value: in.Name,
					}
				},
			},
			{
				Name: "Revision",
				Wide: false,
				Getter: func(in *admin.ConnectorTypeAdminView) dumper.Row {
					// assume we have only the stable channel for now
					if data, ok := in.Channels["stable"].ShardMetadata["connector_revision"]; ok {
						return dumper.Row{
							Value: data.(string),
						}
					}

					return dumper.Row{}
				},
			},
			{
				Name: "Image",
				Wide: false,
				Getter: func(in *admin.ConnectorTypeAdminView) dumper.Row {
					// assume we have only the stable channel for now
					if data, ok := in.Channels["stable"].ShardMetadata["connector_image"]; ok {
						return dumper.Row{
							Value: data.(string),
						}
					}
					if data, ok := in.Channels["stable"].ShardMetadata["container_image"]; ok {
						return dumper.Row{
							Value: data.(string),
						}
					}

					return dumper.Row{}
				},
			},
			{
				Name: "Operator Type",
				Wide: false,
				Getter: func(in *admin.ConnectorTypeAdminView) dumper.Row {
					if data, ok := in.Channels["stable"].ShardMetadata["operators"]; ok {
						operators := data.([]interface{})

						// assume we have a single operator
						return dumper.Row{
							Value: operators[0].(map[string]interface{})["type"].(string),
						}
					}

					return dumper.Row{}
				},
			},
			{
				Name: "Operator Version Range",
				Wide: false,
				Getter: func(in *admin.ConnectorTypeAdminView) dumper.Row {
					if data, ok := in.Channels["stable"].ShardMetadata["operators"]; ok {
						operators := data.([]interface{})

						// assume we have a single operator
						return dumper.Row{
							Value: operators[0].(map[string]interface{})["version"].(string),
						}
					}

					return dumper.Row{}
				},
			},
		},
	}

	dumper.DumpTable(config, out, items.Items)
}
