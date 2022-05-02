package request

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName  = "request"
	CommandAlias = "r"
)

type options struct {
	method string

	f *factory.Factory
}

func NewRequestCommand(f *factory.Factory) *cobra.Command {
	opts := options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     CommandName,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{CommandAlias},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.method != http.MethodGet && opts.method != http.MethodPut && opts.method != http.MethodPost && opts.method != http.MethodPatch && opts.method != http.MethodDelete {
				return fmt.Errorf("unsupported HTTP method: %s, valid values (GET, POST, PUT, PATCH, DELETE)", opts.method)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(&opts, args[0])
		},
	}

	//cmd.Flags().StringVarP(&opts.urlPath, "path", "p", "", "Path to send request. For example /api/connector_mgmt/v1/admin/kafka_connector_clusters")
	cmd.Flags().StringVarP(&opts.method, "method", "X", http.MethodGet, "HTTP method to use. (GET, POST, PUT, PATCH, DELETE)")

	cmd.MarkFlagRequired("path")

	return cmd
}

func run(opts *options, path string) error {
	c, err := service.NewAdminClient(&service.Config{
		F: opts.f,
	})
	if err != nil {
		return err
	}

	var result interface{}
	var httpRes *http.Response

	switch strings.ToUpper(opts.method) {
	case http.MethodGet:
		result, httpRes, err = c.GET(path)
	case http.MethodDelete:
		result, httpRes, err = c.DELETE(path)
	case http.MethodPost:
		result, httpRes, err = c.POST(path, opts.f.IOStreams.In)
	case http.MethodPut:
		result, httpRes, err = c.PUT(path, opts.f.IOStreams.In)
	case http.MethodPatch:
		result, httpRes, err = c.PATCH(path, opts.f.IOStreams.In)
	}

	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return response.Error(err, httpRes)
	}

	if _, err := fmt.Fprint(opts.f.IOStreams.Out, result); err != nil {
		return err
	}

	return nil
}
