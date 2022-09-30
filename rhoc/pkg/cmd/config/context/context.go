package context

import (
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdcontext"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

const (
	CommandName  = "context"
	CommandAlias = "ctx"
)

func NewContextCommand(f *factory.Factory) *cobra.Command {

	// ************************************************
	//
	// ls
	//
	// ************************************************

	ls := cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			items := make([][]string, 0)

			sc, err := f.ServiceContext.Load()
			if err == nil && sc != nil {
				for k, v := range sc.Contexts {
					if k == sc.CurrentContext {
						k = k + " (*)"
					}

					items = append(items, []string{k, v.NamespaceID})
				}
			}

			_, _ = fmt.Fprint(f.IOStreams.Out, "\n")
			_ = dumper.DumpKV(f.IOStreams.Out, []string{"Name", "NamespaceID"}, items)
			_, _ = fmt.Fprint(f.IOStreams.Out, "\n")

			return nil
		},
	}

	// ************************************************
	//
	// use
	//
	// ************************************************

	useOpts := struct {
		name string
	}{}

	use := cobra.Command{
		Use:  "use",
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.ServiceContext.Load()
			if err == nil && sc != nil {
				if _, ok := sc.Contexts[useOpts.name]; !ok {
					return fmt.Errorf("the required context %s does not exists", useOpts.name)
				}

				sc.CurrentContext = useOpts.name

				if err := f.ServiceContext.Save(sc); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmdutil.AddName(&use, &useOpts.name).
		Required().
		Complete(func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var items []string

			sc, err := f.ServiceContext.Load()
			if err == nil && sc != nil {
				items = cmdcontext.List(*sc)
			}

			return items, cobra.ShellCompDirectiveDefault
		})

	// ************************************************
	//
	// create
	//
	// ************************************************

	createOpts := struct {
		name        string
		namespaceId string
		use         bool
	}{}

	create := cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.ServiceContext.Load()
			if err == nil && sc != nil {

				cmdcontext.Create(*sc, createOpts.name)

				c := sc.Contexts[createOpts.name]
				c.NamespaceID = createOpts.namespaceId

				sc.Contexts[createOpts.name] = c

				if createOpts.use {
					sc.CurrentContext = createOpts.name
				}

				if err := f.ServiceContext.Save(sc); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmdutil.AddName(&create, &createOpts.name).Required()
	cmdutil.AddNamespaceID(&create, &createOpts.namespaceId)
	cmdutil.AddUse(&create, &createOpts.use)

	// ************************************************
	//
	// delete
	//
	// ************************************************

	deleteOpts := struct {
		name string
	}{}

	delete := cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.ServiceContext.Load()
			if err == nil && sc != nil {

				delete(sc.Contexts, deleteOpts.name)

				if sc.CurrentContext == deleteOpts.name {
					sc.CurrentContext = ""
				}

				if err := f.ServiceContext.Save(sc); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmdutil.AddName(&delete, &deleteOpts.name).
		Required().
		Complete(func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var items []string

			sc, err := f.ServiceContext.Load()
			if err == nil && sc != nil {
				items = cmdcontext.List(*sc)
			}

			return items, cobra.ShellCompDirectiveDefault
		})

	// ************************************************
	//
	// root
	//
	// ************************************************

	cmd := &cobra.Command{
		Use:     CommandName,
		Aliases: []string{CommandAlias},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(&use)
	cmd.AddCommand(&ls)
	cmd.AddCommand(&create)
	cmd.AddCommand(&delete)

	return cmd
}
