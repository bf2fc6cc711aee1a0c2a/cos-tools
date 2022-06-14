package service

import (
	"net/http"

	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/api"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"

	"golang.org/x/oauth2"
)

type Config struct {
	F *factory.Factory
}

func API(config *Config) (api.API, error) {
	conn, err := config.F.Connection(connection.DefaultConfigRequireMasAuth)
	if err != nil {
		return nil, err
	}

	return conn.API(), nil
}

func NewAdminClient(config *Config) (*AdminAPI, error) {
	a, err := API(config)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: a.GetConfig().MasAccessToken,
		},
	)

	c := admin.NewConfiguration()
	c.Scheme = a.GetConfig().ApiURL.Scheme
	c.Host = a.GetConfig().ApiURL.Host
	c.UserAgent = a.GetConfig().UserAgent
	c.Debug = config.F.Logger.DebugEnabled()
	c.HTTPClient = &http.Client{
		Transport: &oauth2.Transport{
			Base:   a.GetConfig().HTTPClient.Transport,
			Source: oauth2.ReuseTokenSource(nil, ts),
		},
	}

	adminAPI := AdminAPI{
		f:     config.F,
		a:     a,
		c:     c,
		admin: admin.NewAPIClient(c),
	}

	return &adminAPI, nil
}

func Get(f *factory.Factory, format string, getter func(api *AdminAPI) (interface{}, *http.Response, error)) error {
	c, err := NewAdminClient(&Config{
		F: f,
	})
	if err != nil {
		return err
	}

	result, httpRes, err := getter(c)
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return response.Error(err, httpRes)
	}

	return dump.Formatted(f.IOStreams.Out, format, result)
}
