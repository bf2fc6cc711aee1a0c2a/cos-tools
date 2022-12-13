package list

import (
	"io"
	"strconv"
	"time"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/request"
	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/util/duration"
)

func dumpAsTable(out io.Writer, items []deploymentDetail, options request.ListDeploymentsOptions, wide bool, style dumper.TableStyle) {
	config := dumper.TableConfig[deploymentDetail]{
		Style: style,
		Wide:  wide,
		Columns: []dumper.Column[deploymentDetail]{
			{
				Name: "ID",
				Wide: false,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Id,
					}
				},
			},
			{
				Name: "ConnectorID",
				Wide: false,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Spec.ConnectorId,
					}
				},
			},
			{
				Name: "NamespaceID",
				Wide: false,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Spec.NamespaceId,
					}
				},
			},
			{
				Name: "ClusterID",
				Wide: false,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Spec.ClusterId,
					}
				},
			},
			{
				Name: "PlatformId",
				Wide: true,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.PlatformID,
					}
				},
			},
			{
				Name: "Type",
				Wide: true,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Spec.ConnectorTypeId,
					}
				},
			},
			{
				Name: "TypeRevision",
				Wide: !options.ChannelUpdate,
				Getter: func(in *deploymentDetail) dumper.Row {
					if typeRevision, ok := in.Spec.ShardMetadata["connector_revision"]; ok {
						floatRevision, isfloat64 := typeRevision.(float64)
						if isfloat64 {
							return dumper.Row{
								Value: strconv.FormatInt(int64(floatRevision), 10),
							}
						} else {
							return dumper.Row{
								Value: typeRevision.(string),
							}
						}

					}

					return dumper.Row{}
				},
			},
			{
				Name: "UpdatableTypeRevision",
				Wide: !options.ChannelUpdate,
				Getter: func(in *deploymentDetail) dumper.Row {
					var updatableTypeRevision string
					if in.Status.ShardMetadata.Available.Revision == 0 {
						updatableTypeRevision = "-"
					} else {
						updatableTypeRevision = strconv.FormatInt(in.Status.ShardMetadata.Available.Revision, 10)
					}
					return dumper.Row{
						Value: updatableTypeRevision,
					}
				},
			},
			{
				Name: "TypeImage",
				Wide: true,
				Getter: func(in *deploymentDetail) dumper.Row {
					if image, ok := in.Spec.ShardMetadata["connector_image"]; ok {
						return dumper.Row{
							Value: image.(string),
						}
					}
					if image, ok := in.Spec.ShardMetadata["container_image"]; ok {
						return dumper.Row{
							Value: image.(string),
						}
					}

					return dumper.Row{}
				},
			},
			{
				Name: "DesiredState",
				Wide: false,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: string(in.Spec.DesiredState),
					}
				},
			},
			{
				Name: "State",
				Wide: false,
				Getter: func(in *deploymentDetail) dumper.Row {
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
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: strconv.FormatInt(in.Metadata.ResourceVersion, 10),
					}
				},
			},
			{
				Name: "DeploymentVersion",
				Wide: true,
				Getter: func(in *deploymentDetail) dumper.Row {
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
				Getter: func(in *deploymentDetail) dumper.Row {
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
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Metadata.CreatedAt.Format(time.RFC3339),
					}
				},
			},
			{
				Name: "ModifiedAt",
				Wide: true,
				Getter: func(in *deploymentDetail) dumper.Row {
					return dumper.Row{
						Value: in.Metadata.UpdatedAt.Format(time.RFC3339),
					}
				},
			},
		},
	}

	dumper.DumpTable(config, out, items)
}
