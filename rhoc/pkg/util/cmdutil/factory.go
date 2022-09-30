package cmdutil

import (
	"context"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdcontext"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/log"
	"github.com/pkg/errors"
	"github.com/redhat-developer/app-services-cli/pkg/core/config"
	"github.com/redhat-developer/app-services-cli/pkg/core/httputil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/iostreams"
	"github.com/redhat-developer/app-services-cli/pkg/core/localize/goi18n"
	"github.com/redhat-developer/app-services-cli/pkg/core/servicecontext"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/kcconnection"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"net/http"
	"sync"
)

func NewFactory() (*factory.Factory, error) {
	localizer, err := goi18n.New(&goi18n.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to set up localizer")
	}

	var conn connection.Connection
	var connErr error
	var once sync.Once

	f := factory.Factory{
		IOStreams:      iostreams.System(),
		Config:         config.NewFile(),
		Logger:         log.NewLogger(),
		Localizer:      localizer,
		Context:        context.Background(),
		ServiceContext: servicecontext.NewFile(),
	}

	if err := cmdcontext.Init(f.ServiceContext); err != nil {
		return nil, err
	}

	f.Connection = func() (connection.Connection, error) {
		once.Do(func() {
			conn, connErr = NewConnection(&f)
		})

		return conn, connErr
	}

	return &f, nil
}

func NewConnection(f *factory.Factory) (connection.Connection, error) {
	cfg, err := f.Config.Load()
	if err != nil {
		return nil, err
	}

	if cfg.AuthURL == "" {
		cfg.AuthURL = build.AuthURL
	}

	builder := kcconnection.NewConnectionBuilder()

	if cfg.AccessToken != "" {
		builder.WithAccessToken(cfg.AccessToken)
	}
	if cfg.RefreshToken != "" {
		builder.WithRefreshToken(cfg.RefreshToken)
	}
	if cfg.ClientID != "" {
		builder.WithClientID(cfg.ClientID)
	}
	if cfg.Scopes != nil {
		builder.WithScopes(cfg.Scopes...)
	}
	if cfg.APIUrl != "" {
		builder.WithURL(cfg.APIUrl)
	}

	builder.WithAuthURL(cfg.AuthURL)
	builder.WithConsoleURL(build.ConsoleURL)
	builder.WithInsecure(cfg.Insecure)
	builder.WithConfig(f.Config)
	builder.WithLogger(f.Logger)

	builder.WithTransportWrapper(func(a http.RoundTripper) http.RoundTripper {
		return &httputil.LoggingRoundTripper{
			Proxied: a,
			Logger:  f.Logger,
		}
	})

	conn, err := builder.BuildContext(f.Context)
	if err != nil {
		return nil, err
	}

	if err := conn.RefreshTokens(f.Context); err != nil {
		return nil, err
	}

	return conn, nil
}
