package patch

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

type options struct {
	id           string
	skipConfirm  bool
	outputFormat string

	f *factory.Factory
}

func execute(opts options, content map[string]interface{}) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	res, resp, err := c.Clusters().PatchConnector(opts.f.Context, opts.id).Body(content).Execute()
	if err != nil {
		return response.Error(err, resp)
	}
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if resp != nil && resp.StatusCode > 300 {
		return dump.Formatted(opts.f.IOStreams.Out, opts.outputFormat, res)
	}

	return nil
}
