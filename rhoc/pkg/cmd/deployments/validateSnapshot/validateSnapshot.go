package validateSnapshot

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments/snapshot"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/request"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const (
	CommandName  = "validateSnapshot"
	CommandAlias = "vs"
)

type options struct {
	clusterAmount int
	inputFile     string

	f *factory.Factory
}

func NewValidateSnapshotCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Short:   "Validates a previous snapshot against the current state.",
		Use:     CommandName,
		Aliases: []string{CommandAlias},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.clusterAmount < 1 || opts.clusterAmount > 30 {
				return errors.New("clusterAmount must be between 1 and 30")
			}
			if opts.inputFile == "" {
				return errors.New("inputFile must be set")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddClusterAmount(cmd, &opts.clusterAmount)
	cmdutil.AddInputFile(cmd, &opts.inputFile)

	return cmd
}

func run(opts *options) error {
	fmt.Fprintln(opts.f.IOStreams.Out, "Starting snapshot validation")

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

	var snapshots []snapshot.DeploymentSnapshot
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
			snapshots = append(snapshots, snapshot.DeploymentSnapshot{
				ClusterId: cluster.Id,
				Id:        deploy.Id,
				Kind:      deploy.Spec.ConnectorTypeId,
				Status:    deploy.Status.Phase,
			})
		}
	}

	buffer := &bytes.Buffer{}

	fmt.Fprintf(opts.f.IOStreams.Out, "Creating snapshot with %d connector deployments\n", len(snapshots))
	err = dumper.DumpTable(
		dumper.TableConfig[snapshot.DeploymentSnapshot]{
			Style: dumper.TableStyleCSV,
			Wide:  false,
			Columns: []dumper.Column[snapshot.DeploymentSnapshot]{
				{
					Name: "ClusterID",
					Wide: false,
					Getter: func(in *snapshot.DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: in.ClusterId,
						}
					},
				},
				{
					Name: "ID",
					Wide: false,
					Getter: func(in *snapshot.DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: in.Id,
						}
					},
				},
				{
					Name: "ConnectorTypeID",
					Wide: false,
					Getter: func(in *snapshot.DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: in.Kind,
						}
					},
				},
				{
					Name: "Status",
					Wide: false,
					Getter: func(in *snapshot.DeploymentSnapshot) dumper.Row {
						return dumper.Row{
							Value: string(in.Status),
						}
					},
				},
			},
		},
		buffer,
		snapshots,
	)
	if err != nil {
		return err
	}

	return runComparison(opts.f.IOStreams.Out, opts.inputFile, buffer)
}

func runComparison(out io.Writer, inputFile string, newSnapshot io.Reader) error {
	fmt.Fprintln(out, "Comparing old and new snapshots..")
	fmt.Fprintln(out, "--------")

	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	reader2 := csv.NewReader(newSnapshot)
	records2, err := reader2.ReadAll()
	if err != nil {
		return err
	}

	anyDiff := false
	// print unequal lines
	for i := range records {
		diff := false
		for j := range records[i] {
			if records[i][j] != records2[i][j] {
				diff = true
				break
			}
		}
		if diff {
			anyDiff = true
			fmt.Fprintf(out, "Difference located in line %d: \nwas: %v, \nnow: %v\n", i+1, records[i], records2[i])
			fmt.Fprintln(out, "--------")
		}
	}

	if !anyDiff {
		fmt.Fprintln(out, "No diff found!")
	}

	return nil
}
