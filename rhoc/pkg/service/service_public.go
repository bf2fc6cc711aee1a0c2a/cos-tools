package service

import (
	"context"
	"errors"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/public"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/api"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	connectormgmtclient "github.com/redhat-developer/app-services-sdk-go/connectormgmt/apiv1/client"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type PublicAPI struct {
	f      *factory.Factory
	a      api.API
	c      *public.Configuration
	public *public.APIClient
}

func (api *PublicAPI) Context() context.Context {
	return api.f.Context
}
func (api *PublicAPI) ConnectorsApi() connectormgmtclient.ConnectorsApi {
	return api.a.ConnectorsMgmt().ConnectorsApi
}

func (api *PublicAPI) GET(relativePath string) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodGet, "", nil)
}

func (api *PublicAPI) DELETE(relativePath string) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodDelete, "", nil)
}

func (api *PublicAPI) POST(relativePath string, body io.Reader) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodPost, "application/json", body)
}

func (api *PublicAPI) PUT(relativePath string, body io.Reader) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodPut, "application/json", body)
}

func (api *PublicAPI) PATCH(relativePath string, body io.Reader) (interface{}, *http.Response, error) {
	return api.do(relativePath, http.MethodPatch, "application/merge-patch+json", body)
}

func (api *PublicAPI) do(relativePath string, method, contentType string, body io.Reader) (interface{}, *http.Response, error) {
	if !strings.HasPrefix(relativePath, "/api/connector_mgmt/v1/") {
		relativePath = path.Join("/api/connector_mgmt/v1/", relativePath)
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
