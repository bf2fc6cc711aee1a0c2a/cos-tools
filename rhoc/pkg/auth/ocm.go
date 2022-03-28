package auth

import (
	"fmt"

	"github.com/openshift-online/ocm-cli/pkg/ocm"
)

func Token() (string, error) {
	cnx, err := ocm.NewConnection().Build()
	if err != nil {
		return "", fmt.Errorf("Failed to create OCM connection: %v", err)
	}

	defer cnx.Close()

	token, _, err := cnx.Tokens()
	if err != nil {
		return "", fmt.Errorf("Can't get token: %v", err)
	}

	return token, nil
}
