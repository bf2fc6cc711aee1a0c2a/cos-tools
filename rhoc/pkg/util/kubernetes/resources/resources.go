package resources

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

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
