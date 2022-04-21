package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	p "github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

func Add(root *cobra.Command, sub *cobra.Command) {
	if err := bindPFlags(sub); err != nil {
		panic(err)
	}

	root.AddCommand(sub)

}

func Bind(root *cobra.Command, subs ...*cobra.Command) {
	if err := bindPFlags(root); err != nil {
		panic(err)
	}

	for _, s := range subs {
		Add(root, s)
	}
}

func bindPFlags(cmd *cobra.Command) (err error) {
	pl := p.NewClient()

	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if err != nil {
			return
		}

		err = bindFlag(pl, flag)
	})

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if err != nil {
			return
		}

		err = bindFlag(pl, flag)
	})

	return err
}

func bindFlag(pl *p.Client, flag *pflag.Flag) error {
	name := flag.Name
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, ".", "-")

	if err := viper.BindPFlag(name, flag); err != nil {
		return fmt.Errorf("error binding flag %s to viper: %v", flag.Name, err)
	}

	// this is a little bit of an hack to register plural version of properties
	// based on the naming conventions used by the flag type because it is not
	// possible to know what is the type of a flag
	flagType := strings.ToUpper(flag.Value.Type())
	if strings.Contains(flagType, "SLICE") || strings.Contains(flagType, "ARRAY") {
		if err := viper.BindPFlag(pl.Plural(name), flag); err != nil {
			return fmt.Errorf("error binding plural flag %s to viper: %v", flag.Name, err)
		}
	}

	return nil
}

func PromptConfirm(format string, args ...interface{}) (bool, error) {
	promptConfirm := survey.Confirm{
		Message: fmt.Sprintf(format, args...),
	}

	var confirmDelete bool
	if err := survey.AskOne(&promptConfirm, &confirmDelete); err != nil {
		return false, err
	}

	return confirmDelete, nil
}
