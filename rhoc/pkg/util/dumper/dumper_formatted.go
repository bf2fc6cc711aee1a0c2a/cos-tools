package dumper

import (
	"io"
	"net/http"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
)

func NewFormatted(format string, httpRes *http.Response, httpErr error) Formatted {
	return Formatted{
		format:       format,
		httpResponse: httpRes,
		httpError:    httpErr,
	}
}

type Formatted struct {
	format       string
	httpResponse *http.Response
	httpError    error
}

func (f Formatted) Dump(out io.Writer, data interface{}) error {
	if f.httpError != nil {
		return response.Error(f.httpError, f.httpResponse)
	}

	return dump.Formatted(out, f.format, data)
}
