package service

import (
	"strconv"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/request"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/response"
)

func ListClusters(c *AdminAPI, opts request.ListOptions) (admin.ConnectorClusterAdminList, error) {
	items := admin.ConnectorClusterAdminList{
		Kind:  "ConnectorClusterAdminList",
		Items: make([]admin.ConnectorClusterAdminView, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Clusters().ListConnectorClusters(c.Context())
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}
		if opts.Search != "" {
			e = e.Search(opts.Search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListNamespacesForCluster(c *AdminAPI, opts request.ListOptions, clusterID string) (admin.ConnectorNamespaceList, error) {
	items := admin.ConnectorNamespaceList{
		Kind:  "ConnectorClusterList",
		Items: make([]admin.ConnectorNamespace, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Clusters().GetClusterNamespaces(c.Context(), clusterID)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}
		if opts.Search != "" {
			e = e.Search(opts.Search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListNamespaces(c *AdminAPI, opts request.ListOptions) (admin.ConnectorNamespaceList, error) {
	items := admin.ConnectorNamespaceList{
		Kind:  "ConnectorClusterList",
		Items: make([]admin.ConnectorNamespace, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Namespaces().GetConnectorNamespaces(c.Context())
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}
		if opts.Search != "" {
			e = e.Search(opts.Search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListConnectorsForCluster(c *AdminAPI, opts request.ListOptions, clusterID string) (admin.ConnectorAdminViewList, error) {
	items := admin.ConnectorAdminViewList{
		Kind:  "ConnectorAdminViewList",
		Items: make([]admin.ConnectorAdminView, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Clusters().GetClusterConnectors(c.Context(), clusterID)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}
		if opts.Search != "" {
			e = e.Search(opts.Search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if result == nil || len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListConnectorsForNamespace(c *AdminAPI, opts request.ListOptions, namespaceID string) (admin.ConnectorAdminViewList, error) {
	items := admin.ConnectorAdminViewList{
		Kind:  "ConnectorAdminViewList",
		Items: make([]admin.ConnectorAdminView, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Clusters().GetNamespaceConnectors(c.Context(), namespaceID)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}
		if opts.Search != "" {
			e = e.Search(opts.Search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if result == nil || len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListDeploymentsForCluster(c *AdminAPI, opts request.ListDeploymentsOptions, clusterID string) (admin.ConnectorDeploymentAdminViewList, error) {
	items := admin.ConnectorDeploymentAdminViewList{
		Kind:  "ConnectorDeploymentAdminViewList",
		Items: make([]admin.ConnectorDeploymentAdminView, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Clusters().GetClusterDeployments(c.Context(), clusterID)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}

		if opts.DanglingDeployments {
			e = e.DanglingDeployments(opts.DanglingDeployments)
		}

		if opts.ChannelUpdate {
			e = e.ChannelUpdates(opts.ChannelUpdate)
		}

		if opts.OperatorUpdate {
			e = e.OperatorUpdates(opts.OperatorUpdate)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListDeploymentsForNamespace(c *AdminAPI, opts request.ListOptions, namespaceId string) (admin.ConnectorDeploymentAdminViewList, error) {
	items := admin.ConnectorDeploymentAdminViewList{
		Kind:  "ConnectorDeploymentAdminViewList",
		Items: make([]admin.ConnectorDeploymentAdminView, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Clusters().GetNamespaceDeployments(c.Context(), namespaceId)
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func ListConnectorTypes(c *AdminAPI, opts request.ListOptions) (admin.ConnectorTypeAdminViewList, error) {
	items := admin.ConnectorTypeAdminViewList{
		Kind:  "ConnectorTypeAdminViewList",
		Items: make([]admin.ConnectorTypeAdminView, 0),
		Total: 0,
		Size:  0,
	}

	for i := opts.Page; i == opts.Page || opts.AllPages; i++ {
		e := c.Catalog().GetConnectorTypes(c.Context())
		e = e.Page(strconv.Itoa(i))
		e = e.Size(strconv.Itoa(opts.Limit))

		if opts.OrderBy != "" {
			e = e.OrderBy(opts.OrderBy)
		}
		if opts.Search != "" {
			e = e.Search(opts.Search)
		}

		result, httpRes, err := e.Execute()

		if httpRes != nil {
			defer func() {
				_ = httpRes.Body.Close()
			}()
		}
		if err != nil {
			items.Size = 0
			items.Total = 0
			items.Items = nil

			return items, response.Error(err, httpRes)
		}
		if len(result.Items) == 0 {
			break
		}

		items.Items = append(items.Items, result.Items...)
		items.Size = int32(len(items.Items))
		items.Total = result.Total
	}

	return items, nil
}

func GetClusterForNamespace(c *AdminAPI, namespaceID string) (*admin.ConnectorClusterAdminView, error) {
	result, httpRes, err := c.Clusters().GetConnectorNamespace(c.Context(), namespaceID).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return nil, response.Error(err, httpRes)
	}

	return GetClusterByID(c, result.ClusterId)
}

func GetClusterByID(c *AdminAPI, clusterID string) (*admin.ConnectorClusterAdminView, error) {
	result, httpRes, err := c.Clusters().GetConnectorCluster(c.Context(), clusterID).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return nil, response.Error(err, httpRes)
	}

	return result, nil
}

func GetConnectorByID(c *AdminAPI, connectorID string) (*admin.ConnectorAdminView, error) {
	result, httpRes, err := c.Clusters().GetConnector(c.Context(), connectorID).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return nil, response.Error(err, httpRes)
	}

	return result, nil
}

func GetNamespaceByID(c *AdminAPI, namespaceID string) (*admin.ConnectorNamespace, error) {
	result, httpRes, err := c.Clusters().GetConnectorNamespace(c.Context(), namespaceID).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return nil, response.Error(err, httpRes)
	}

	return result, nil
}

func GetDeploymentByID(c *AdminAPI, clusterID string, deploymentID string) (*admin.ConnectorDeploymentAdminView, error) {
	result, httpRes, err := c.Clusters().GetConnectorDeployment(c.Context(), clusterID, deploymentID).Execute()
	if httpRes != nil {
		defer func() {
			_ = httpRes.Body.Close()
		}()
	}
	if err != nil {
		return nil, response.Error(err, httpRes)
	}

	return result, nil
}
