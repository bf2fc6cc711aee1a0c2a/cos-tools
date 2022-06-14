package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/api"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

type AdminAPI struct {
	f     *factory.Factory
	a     api.API
	c     *admin.Configuration
	admin *admin.APIClient
}

func (api *AdminAPI) Context() context.Context {
	return api.f.Context
}

func (api *AdminAPI) Clusters() *admin.ConnectorClustersAdminApiService {
	return api.admin.ConnectorClustersAdminApi
}

func (api *AdminAPI) Namespaces() *admin.ConnectorNamespacesAdminApiService {
	return api.admin.ConnectorNamespacesAdminApi
}

func (api *AdminAPI) Catalog() *admin.ConnectorTypesApiService {
	return api.admin.ConnectorTypesApi
}

func (api *AdminAPI) GET(relativePath string) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodGet, "", nil)
}

func (api *AdminAPI) DELETE(relativePath string) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodDelete, "", nil)
}

func (api *AdminAPI) POST(relativePath string, body io.Reader) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodPost, "application/json", body)
}

func (api *AdminAPI) PUT(relativePath string, body io.Reader) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodPut, "application/json", body)
}

func (api *AdminAPI) PATCH(relativePath string, body io.Reader) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodPatch, "application/merge-patch+json", body)
}

func (api *AdminAPI) do(relativePath string, method, contentType string, body io.Reader) (interface{}, *http.Response, error) {
	if !strings.HasPrefix(relativePath, "/api/connector_mgmt/v1/admin/") {
		relativePath = path.Join("/api/connector_mgmt/v1/admin/", relativePath)
	}

	u, err := url.Parse(relativePath)
	if err != nil {
		return nil, nil, err
	}

	url := *api.a.GetConfig().ApiURL
	url.Path = u.Path
	url.RawQuery = u.RawQuery

	var req *http.Request

	if body != nil {
		r, w := io.Pipe()

		go func() {
			defer w.Close()

			if _, err := io.Copy(w, body); err != nil {
				w.CloseWithError(err)
				return
			}
			if err := w.Close(); err != nil {
				w.CloseWithError(err)
				return
			}
		}()

		req, err = http.NewRequest(method, url.String(), r)
		if err != nil {
			return nil, nil, err
		}

		req = req.WithContext(api.f.Context)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Accept", "application/json")
	} else {
		req, err = http.NewRequest(method, url.String(), nil)
		if err != nil {
			return err, nil, err
		}

		req = req.WithContext(api.f.Context)
		req.Header.Set("Accept", "application/json")
	}

	resp, err := api.c.HTTPClient.Do(req)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode > http.StatusBadRequest {
		return nil, resp, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	return string(b), resp, err
}
