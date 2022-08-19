package resources

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/collections"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
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
