package resources

import (
	"context"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/collections"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"strings"
)

func Discover(client discovery.DiscoveryInterface) ([]schema.GroupVersionResource, error) {
	resources, err := client.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	answer := make([]schema.GroupVersionResource, 0)

	for _, group := range resources {
		gv, err := schema.ParseGroupVersion(group.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, r := range group.APIResources {
			if !r.Namespaced {
				continue
			}
			if !collections.Contains(r.Verbs, "list") {
				continue
			}

			answer = append(answer, schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: r.Name,
			})
		}
	}

	return answer, nil
}

func Parse(resources []string) ([]schema.GroupVersionResource, error) {
	answer := make([]schema.GroupVersionResource, 0)
	for _, r := range resources {
		s := strings.Split(r, ":")
		if len(s) != 2 {
			continue
		}

		gv, err := schema.ParseGroupVersion(s[0])
		if err != nil {
			return nil, err
		}

		answer = append(answer, schema.GroupVersionResource{
			Group:    gv.Group,
			Version:  gv.Version,
			Resource: s[1],
		})
	}

	return answer, nil
}

func List(ctx context.Context, client dynamic.Interface, resources []schema.GroupVersionResource, connectorId string) ([]unstructured.Unstructured, error) {
	items := make([]unstructured.Unstructured, 0)

	for _, gvr := range resources {

		if gvr.Resource == "" {
			continue
		}

		resList, err := client.Resource(gvr).List(ctx, metav1.ListOptions{
			LabelSelector: "cos.bf2.org/connector.id=" + connectorId,
		})
		if err != nil {
			return nil, err
		}

		//case item.GetAPIVersion() == "v1" && item.GetKind() == "Secret":
		//continue

		for i := range resList.Items {
			// remove managed fields as they are only making noise
			resList.Items[i].SetManagedFields(nil)

			items = append(items, resList.Items[i])
		}
	}

	return items, nil
}
