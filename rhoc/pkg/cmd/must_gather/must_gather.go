package must_gather

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes/pods"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes/resources"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"k8s.io/apimachinery/pkg/runtime/schema"
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
	file      string
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

			for i := range opts.gvrs {
				// exclude secrets
				if opts.gvrs[i].Group == "" && opts.gvrs[i].Version == "v1" && opts.gvrs[i].Resource == "secrets" {
					opts.gvrs[i].Group = ""
					opts.gvrs[i].Version = ""
					opts.gvrs[i].Resource = ""
				}
			}

			var o io.Writer
			var mustClose bool

			if opts.file == "" {
				o = f.IOStreams.Out
				mustClose = false
			} else {
				f, err := os.Create(opts.file)
				if err != nil {
					return err
				}

				o = f
				mustClose = true
			}

			out := bufio.NewWriter(o)
			defer func() {
				_ = out.Flush()

				if mustClose {
					if c, ok := o.(io.Closer); ok {
						_ = c.Close()
					}
				}
			}()

			items, err := resources.List(f.Context, dynamicClient, opts.gvrs, opts.id)
			if err != nil {
				return err
			}

			if len(items) != 0 {
				raw, err := yaml.Marshal(items)
				if err != nil {
					return err
				}

				_, err = out.Write(raw)
				if err != nil {
					return err
				}
			}

			if opts.logs {
				for _, item := range items {
					if item.GetAPIVersion() == "v1" && item.GetKind() == "Pod" {
						containers, err := pods.ListContainers(f.Context, client, item.GetNamespace(), item.GetName())
						if err != nil {
							return err
						}

						for _, container := range containers {
							_, err = fmt.Fprintf(
								out,
								"%s/%s:%s:%s@%s\n",
								item.GetAPIVersion(),
								item.GetKind(),
								item.GetNamespace(),
								item.GetName(),
								container)

							if err != nil {
								return err
							}

							err := pods.Logs(f.Context, client, item.GetNamespace(), item.GetName(), container, out)
							if err != nil && !errors.Is(err, io.EOF) {
								return err
							}
						}

						_, err = out.Write([]byte{'\n'})
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
	}

	cmdutil.AddID(cmd, &opts.id).Required()
	cmdutil.AddFile(cmd, &opts.file)

	cmd.Flags().BoolVar(&opts.logs, "logs", opts.logs, "Include logs")
	cmd.Flags().StringSliceVar(&opts.resources, "resource", nil, "resources to include apiVersion:plural")

	configFlags.AddFlags(cmd.PersistentFlags())

	return cmd
}
