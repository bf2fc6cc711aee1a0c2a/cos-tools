package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
)

func Error(err error, resp *http.Response) error {
	if resp != nil && resp.StatusCode == http.StatusInternalServerError {
		e, _ := service.ReadError(resp)
		if e.Reason != "" {
			err = fmt.Errorf("%s: [%w]", err.Error(), errors.New(e.Reason))
		}
	}
	return err
}
