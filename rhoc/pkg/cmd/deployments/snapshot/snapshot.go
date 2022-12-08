package snapshot

import (
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/request"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"time"
)

const (
	CommandName  = "snapshot"
	CommandAlias = "ss"
)

type options struct {
	clusterAmount int
	outputFile    string

	f *factory.Factory
}

func NewSnapshotCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Short:   "Takes a snapshot of deployments in oldest and ready clusters.",
		Use:     CommandName,
		Aliases: []string{CommandAlias},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.clusterAmount < 1 || opts.clusterAmount > 30 {
				return errors.New("clusterAmount must be between 1 and 30")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddClusterAmount(cmd, &opts.clusterAmount)
	cmdutil.AddOutputFile(cmd, &opts.outputFile)

	return cmd
}

func run(opts *options) error {
	fmt.Fprintln(opts.f.IOStreams.Out, "Starting snapshot creation")

	if opts.outputFile == "" {
		defaultFilename := fmt.Sprintf("snapshot%s.csv", time.Now().Format("20060102150405"))
		fmt.Fprintf(opts.f.IOStreams.Out, "outputFile wasn't set, results will be written to %s\n", defaultFilename)
		opts.outputFile = defaultFilename
	}

	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.f.IOStreams.Out, "Querying for clusters..")

	clusterListOpts := request.ListOptions{
		Page:     0,
		Limit:    opts.clusterAmount,
		AllPages: false,
		OrderBy:  "created_at",
		Search:   "state=ready",
	}
	clusters, err := service.ListClusters(c, clusterListOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.f.IOStreams.Out, "Got %d clusters\n", clusters.Size)

	var snapshots []DeploymentSnapshot
	for _, cluster := range clusters.Items {
		fmt.Fprintf(opts.f.IOStreams.Out, "Querying for deployments in cluster %s..\n", cluster.Id)

		deployListOpts := request.ListDeploymentsOptions{
			ListOptions: request.ListOptions{
				AllPages: false,
				Limit:    30,
				Page:     0,
			},
			ChannelUpdate:       false,
			DanglingDeployments: false,
		}

		deployments, err := service.ListDeploymentsForCluster(c, deployListOpts, cluster.Id)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.f.IOStreams.Out, "Got %d deployments in cluster %s\n", deployments.Size, cluster.Id)

		for _, deploy := range deployments.Items {
			snapshots = append(snapshots, DeploymentSnapshot{
				ClusterId: cluster.Id,
				Id:        deploy.Id,
				Kind:      deploy.Spec.ConnectorTypeId,
				Status:    deploy.Status.Phase,
			})
		}
	}

	outputWriter, err := cmdutil.NewOutputFileWriter(opts.outputFile)
	if err != nil {
		return err
	}
	defer outputWriter.Close()

	fmt.Fprintf(opts.f.IOStreams.Out, "Creating snapshot with %d connector deployments\n", len(snapshots))
	err = dumper.DumpTable(
		dumper.TableConfig[DeploymentSnapshot]{
			Style: dumper.TableStyleCSV,
			Wide:  false,
			Columns: []dumper.Column[DeploymentSnapshot]{
				{
					Name: "ClusterID",
					Wide: false,
					Getter: func(in *DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: in.ClusterId,
						}
					},
				},
				{
					Name: "ID",
					Wide: false,
					Getter: func(in *DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: in.Id,
						}
					},
				},
				{
					Name: "ConnectorTypeID",
					Wide: false,
					Getter: func(in *DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: in.Kind,
						}
					},
				},
				{
					Name: "Status",
					Wide: false,
					Getter: func(in *DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: string(in.Status),
						}
					},
				},
			},
		},
		outputWriter,
		snapshots,
	)
	if err != nil {
		return err
	}
	fmt.Fprintf(opts.f.IOStreams.Out, "Snapshot created with %d connector deployments\n", len(snapshots))

	return nil
}
