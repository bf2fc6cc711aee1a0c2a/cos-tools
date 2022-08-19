package must_gather

import (
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes/pods"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes/resources"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

type options struct {
	id        string
	logs      bool
	resources []string
	gvrs      []schema.GroupVersionResource
	f         *factory.Factory
}

func NewMustGatherCommand(f *factory.Factory) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)
	_ = printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.YAMLPrinter{})

	opts := options{
		f:    f,
		logs: false,
	}

	cmd := &cobra.Command{
		Use:     "must-gather",
		Aliases: []string{"mg"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := configFlags.ToRESTConfig()
			if err != nil {
				return err
			}
			client, _ := kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}
			dynamicClient, err := dynamic.NewForConfig(config)
			if err != nil {
				return err
			}

			if opts.resources == nil {
				discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
				if err != nil {
					return err
				}

				opts.gvrs, err = resources.Discover(discoveryClient)
				if err != nil {
					return err
				}
			} else {
				opts.gvrs, err = resources.Parse(opts.resources)
				if err != nil {
					return err
				}
			}

			ulist := make([]map[string]interface{}, 0)

			for _, gvr := range opts.gvrs {

				resList, err := dynamicClient.Resource(gvr).List(f.Context, metav1.ListOptions{
					LabelSelector: "cos.bf2.org/connector.id=" + opts.id,
				})
				if err != nil {
					return err
				}

				for _, res := range resList.Items {
					switch {
					case res.GetAPIVersion() == "v1" && res.GetKind() == "Secret":
						continue
					case res.GetAPIVersion() == "v1" && res.GetKind() == "ConfigMap":
						continue
					default:
						fmt.Printf("Gathering -> %s:%s\n", res.GetAPIVersion(), gvr.Resource)

						// remove managed fields as they are only making noise
						res.SetManagedFields(nil)

						ulist = append(ulist, res.Object)

						if res.GetAPIVersion() == "v1" && res.GetKind() == "Pod" && opts.logs {
							containers, err := pods.ListContainers(f.Context, client, res.GetNamespace(), res.GetName())
							if err != nil {
								return err
							}

							for _, container := range containers {
								err := pods.Logs(f.Context, client, res.GetNamespace(), res.GetName(), container, os.Stdout)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}

			if len(ulist) != 0 {
				out, err := yaml.Marshal(ulist)
				if err != nil {
					return err
				}

				fmt.Println(string(out))
			}

			return nil
		},
	}

	cmdutil.AddID(cmd, &opts.id).Required()

	cmd.Flags().BoolVar(&opts.logs, "logs", opts.logs, "Include logs")
	cmd.Flags().StringSliceVar(&opts.resources, "resource", nil, "resources to include apiVersion:plural")

	configFlags.AddFlags(cmd.PersistentFlags())

	return cmd
}
