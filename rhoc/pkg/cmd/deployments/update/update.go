package update

import (
	"errors"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName = "update"
)

type options struct {
	id           string
	clusterID    string
	outputFormat string
	operatorId   string
	revision     int64

	f *factory.Factory
}

func NewUpdateCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:  CommandName,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.ValidateOutputs(cmd); err != nil {
				return err
			}
			// validate that at least one between revision and operator-id is specified
			revision, err := cmd.Flags().GetInt64("revision")
			if err != nil {
				return err
			}
			operatorId, err := cmd.Flags().GetString("operator-id")
			if err != nil {
				return err
			}
			if revision == 0 && operatorId == "" {
				return errors.New("at least one between revision and operator-id must be specified (specifying both is also valid)")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts)
		},
	}

	cmdutil.AddOutput(cmd, &opts.outputFormat)
	cmdutil.AddID(cmd, &opts.id).Required()
	cmdutil.AddClusterID(cmd, &opts.clusterID).Required()

	cmdutil.AddRevision(cmd, &opts.revision)
	cmdutil.AddOperatorId(cmd, &opts.operatorId)

	return cmd
}

func run(opts *options) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	updateBody := map[string]interface{}{
		"spec": map[string]interface{}{},
	}

	if opts.revision != 0 {
		updateBody["spec"].(map[string]interface{})["shard_metadata"] = map[string]interface{}{
			"connector_revision": opts.revision,
		}
	}

	if opts.operatorId != "" {
		updateBody["spec"].(map[string]interface{})["operator_id"] = opts.operatorId
	}

	result, httpRes, err := c.Clusters().PatchConnectorClusterDeploymentAdmin(opts.f.Context, opts.clusterID, opts.id).Body(updateBody).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return response.Error(err, httpRes)
	}

	return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, result)
}
