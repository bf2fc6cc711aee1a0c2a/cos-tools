package ocm

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build"
	ocmcfg "github.com/openshift-online/ocm-cli/pkg/config"
	ocmcli "github.com/openshift-online/ocm-cli/pkg/ocm"
	ocmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/pkg/errors"
)

func Cluster(externalID string) (*ocmv1.Cluster, error) {
	if externalID == "" {
		return nil, nil
	}

	_, err := ocmcfg.Load()
	if err != nil {
		return nil, err
	}

	// Create the client for the OCM API:
	connection, err := ocmcli.NewConnection().Build()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = connection.Close()
	}()

	search := "external_id='" + externalID + "'"

	request := connection.ClustersMgmt().V1().Clusters().List().Search(search)
	request.Size(build.DefaultPageSize)
	request.Page(build.DefaultPageNumber)

	response, err := request.Send()
	if err != nil {
		return nil, errors.Wrapf(err, "can't retrieve cluster %s", externalID)
	}
	if response.Size() != 1 {
		return nil, errors.Wrapf(err, "can't retrieve cluster %s, found %d", externalID, response.Size())
	}

	return response.Items().Get(0), nil
}
