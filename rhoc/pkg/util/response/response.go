package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
)

func ReadError(response *http.Response) (admin.Error, error) {
	serviceError := admin.Error{}
	err := json.NewDecoder(response.Body).Decode(&serviceError)

	return serviceError, err
}

func Error(err error, resp *http.Response) error {
	if resp != nil && resp.StatusCode == http.StatusInternalServerError {
		e, _ := ReadError(resp)
		if e.Reason != "" {
			err = fmt.Errorf("%s: [%w]", err.Error(), errors.New(e.Reason))
		}
	}
	return err
}
