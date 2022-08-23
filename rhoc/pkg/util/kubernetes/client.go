package kubernetes

import (
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/collections"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes/pods"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"io"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	d "k8s.io/client-go/dynamic"
	k "k8s.io/client-go/kubernetes"
	r "k8s.io/client-go/rest"
)

type Client struct {
	f      *factory.Factory
	config *r.Config
	C      k.Interface
}

func NewClient(f *factory.Factory, config *r.Config) (*Client, error) {
	client, err := k.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	answer := Client{
		f:      f,
		config: config,
		C:      client,
	}

	return &answer, nil
}

func (in *Client) ServerResources() ([]schema.GroupVersionResource, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(in.config)
	if err != nil {
		return nil, err
	}

	resources, err := dc.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	answer := make([]schema.GroupVersionResource, 0)

	for _, group := range resources {
		gv, err := schema.ParseGroupVersion(group.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, res := range group.APIResources {
			if !res.Namespaced {
				continue
			}
			if !collections.Contains(res.Verbs, "list") {
				continue
			}

			answer = append(answer, schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: res.Name,
			})
		}
	}

	return answer, nil
}

func (in *Client) List(resources []schema.GroupVersionResource, options metav1.ListOptions) ([]unstructured.Unstructured, error) {
	dc, err := d.NewForConfig(in.config)
	if err != nil {
		return nil, err
	}

	items := make([]unstructured.Unstructured, 0)

	for _, gvr := range resources {
		if gvr.Resource == "" {
			continue
		}

		resList, err := dc.Resource(gvr).List(in.f.Context, options)
		if err != nil {
			switch {
			case kerr.IsUnauthorized(err):
				in.f.Logger.Debugf("unauthorized access to %s:%s", gvr.GroupVersion(), gvr.Resource)
				continue
			case kerr.IsForbidden(err):
				in.f.Logger.Debugf("forbidden access to %s:%s", gvr.GroupVersion(), gvr.Resource)
				continue
			default:
				in.f.Logger.Debugf("error listing %s:%s, %s", gvr.GroupVersion(), gvr.Resource, err.Error())
				return nil, err
			}
		}

		in.f.Logger.Debugf("found %d %s:%s resources", len(resList.Items), gvr.GroupVersion(), gvr.Resource)

		for i := range resList.Items {
			// remove managed fields as they are only making noise
			resList.Items[i].SetManagedFields(nil)

			items = append(items, resList.Items[i])
		}
	}

	return items, nil
}

func (in *Client) Logs(namespace string, name string, writer io.Writer) error {
	containers, err := pods.ListContainers(in.f.Context, in.C, namespace, name)
	if err != nil {
		return err
	}

	for _, container := range containers {
		err := in.LogsForContainer(namespace, name, container, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *Client) LogsForContainer(namespace string, name string, container string, writer io.Writer) error {
	_, err := fmt.Fprintf(
		writer,
		"v1/pods:%s:%s@%s\n",
		namespace,
		name,
		container)

	if err != nil {
		return err
	}

	err = pods.Logs(in.f.Context, in.C, namespace, name, container, writer)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
