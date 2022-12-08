package cmdutil

import (
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	p "github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	OutputFormatWide = "wide"
	OutputFormatCSV  = "csv"
)

func Add(root *cobra.Command, sub *cobra.Command) *cobra.Command {
	if err := bindPFlags(sub); err != nil {
		panic(err)
	}

	root.AddCommand(sub)

	return root
}

func Bind(root *cobra.Command, subs ...*cobra.Command) *cobra.Command {
	if err := bindPFlags(root); err != nil {
		panic(err)
	}

	for _, s := range subs {
		Add(root, s)
	}

	return root
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

func ValidOutputs() []string {
	validVals := make([]string, 0, len(flagutil.ValidOutputFormats)+1)
	validVals = append(validVals, flagutil.ValidOutputFormats...)
	validVals = append(validVals, OutputFormatWide)
	validVals = append(validVals, OutputFormatCSV)

	return validVals
}

func ValidateOutputs(cmd *cobra.Command) error {
	formats, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	if formats != "" && !flagutil.IsValidInput(formats, ValidOutputs()...) {
		return flagutil.InvalidValueError("output", formats, ValidOutputs()...)
	}

	return nil
}

func AddOutput(cmd *cobra.Command, output *string) *FlagOptions {
	validVals := ValidOutputs()
	name := "output"

	cmd.Flags().StringVarP(
		output,
		name,
		"o",
		dump.EmptyFormat,
		"Specify the output format. Choose from: "+strings.Join(validVals, ", "),
	)

	_ = cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return validVals, cobra.ShellCompDirectiveNoSpace
	})

	return withFlagOptions(cmd, name)
}

func AddPage(cmd *cobra.Command, output *int) *FlagOptions {
	name := "page"

	cmd.Flags().IntVarP(
		output,
		name,
		"p",
		build.DefaultPageNumber,
		"Page index",
	)

	return withFlagOptions(cmd, name)
}

func AddLimit(cmd *cobra.Command, output *int) *FlagOptions {
	name := "limit"

	cmd.Flags().IntVarP(
		output,
		name,
		"l",
		build.DefaultPageSize,
		"Number of items in each page",
	)

	return withFlagOptions(cmd, name)
}

func AddAllPages(cmd *cobra.Command, output *bool) *FlagOptions {
	name := "all-pages"

	cmd.Flags().BoolVar(
		output,
		name,
		false,
		"Grab all pages",
	)

	return withFlagOptions(cmd, name)
}

func AddOrderBy(cmd *cobra.Command, output *string) *FlagOptions {
	name := "order-by"

	cmd.Flags().StringVar(
		output,
		name,
		*output,
		"Specifies the order by criteria",
	)

	return withFlagOptions(cmd, name)
}

func AddSearch(cmd *cobra.Command, output *string) *FlagOptions {
	name := "search"

	cmd.Flags().StringVar(
		output,
		name,
		*output,
		"Search criteria",
	)

	return withFlagOptions(cmd, name)
}

func AddClusterID(cmd *cobra.Command, output *string) *FlagOptions {
	name := "cluster-id"

	cmd.Flags().StringVarP(
		output,
		name,
		"c",
		*output,
		"Cluster ID",
	)

	return withFlagOptions(cmd, name)
}

func AddNamespaceID(cmd *cobra.Command, output *string) *FlagOptions {
	name := "namespace-id"

	cmd.Flags().StringVarP(
		output,
		name,
		"n",
		*output,
		"Namespace ID",
	)

	return withFlagOptions(cmd, name)
}

func AddID(cmd *cobra.Command, output *string) *FlagOptions {
	name := "id"

	cmd.Flags().StringVar(
		output,
		name,
		*output,
		"ID",
	)

	return withFlagOptions(cmd, name)
}

func AddTenantKind(cmd *cobra.Command, output *string) *FlagOptions {
	name := "tenant-kind"

	cmd.Flags().StringVar(
		output,
		name,
		*output,
		"Tenant Kind",
	)

	return withFlagOptions(cmd, name)
}

func AddTenantID(cmd *cobra.Command, output *string) *FlagOptions {
	name := "tenant-id"

	cmd.Flags().StringVar(
		output,
		name,
		*output,
		"Tenant ID",
	)

	return withFlagOptions(cmd, name)
}

func AddName(cmd *cobra.Command, output *string) *FlagOptions {
	name := "name"

	cmd.Flags().StringVar(
		output,
		name,
		*output,
		"Name",
	)

	return withFlagOptions(cmd, name)
}

func AddForce(cmd *cobra.Command, output *bool) *FlagOptions {
	name := "force"

	cmd.Flags().BoolVarP(
		output,
		name,
		"f",
		false,
		"Force",
	)

	return withFlagOptions(cmd, name)
}

func AddUse(cmd *cobra.Command, output *bool) *FlagOptions {
	name := "use"

	cmd.Flags().BoolVar(
		output,
		name,
		false,
		"Use",
	)

	return withFlagOptions(cmd, name)
}

func AddYes(cmd *cobra.Command, yes *bool) *FlagOptions {
	name := "yes"

	cmd.Flags().BoolVarP(
		yes,
		name,
		"y",
		false,
		"Skip confirmation of this action",
	)

	return withFlagOptions(cmd, name)
}

func AddChannelUpdate(cmd *cobra.Command, output *bool) *FlagOptions {
	name := "channel-update"

	cmd.Flags().BoolVar(
		output,
		name,
		false,
		"Filter deployment with channel updates available",
	)

	return withFlagOptions(cmd, name)
}

func AddDanglingDeployments(cmd *cobra.Command, output *bool) *FlagOptions {
	name := "dangling-deployments"

	cmd.Flags().BoolVar(
		output,
		name,
		false,
		"Filter not deleted deployment referring to a deleted connector",
	)

	return withFlagOptions(cmd, name)
}

func AddRevision(cmd *cobra.Command, output *int64) *FlagOptions {
	name := "revision"

	cmd.Flags().Int64Var(
		output,
		name,
		0,
		"Revision to update to",
	)

	return withFlagOptions(cmd, name)
}

func AddClusterAmount(cmd *cobra.Command, output *int) *FlagOptions {
	name := "clusterAmount"

	cmd.Flags().IntVar(
		output,
		name,
		3,
		"Amount of clusters to query when taking the snapshot",
	)

	return withFlagOptions(cmd, name)
}

func AddOutputFile(cmd *cobra.Command, output *string) *FlagOptions {
	name := "outputFile"

	cmd.Flags().StringVar(
		output,
		name,
		"",
		"File to write the output to",
	)

	return withFlagOptions(cmd, name)
}

func AddInputFile(cmd *cobra.Command, output *string) *FlagOptions {
	name := "inputFile"

	cmd.Flags().StringVar(
		output,
		name,
		"",
		"File to read the input from",
	)

	return withFlagOptions(cmd, name)
}

func AddFile(cmd *cobra.Command, output *string) *FlagOptions {
	name := "file"

	cmd.Flags().StringVar(
		output,
		name,
		"",
		"file",
	)

	return withFlagOptions(cmd, name)
}

func withFlagOptions(cmd *cobra.Command, flagName string) *FlagOptions {
	options := FlagOptions{}

	options.Required = func() *FlagOptions {
		_ = cmd.MarkFlagRequired(flagName)
		return &options
	}
	options.Complete = func(f func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)) *FlagOptions {
		_ = cmd.RegisterFlagCompletionFunc(flagName, f)
		return &options
	}

	return &options
}

// FlagOptions defines additional flag options
type FlagOptions struct {
	Required func() *FlagOptions
	Complete func(func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)) *FlagOptions
}
