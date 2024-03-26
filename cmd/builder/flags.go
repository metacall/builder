package builder

import (
	"github.com/spf13/cobra"
)

type DepsImageFlags struct {
	Branch string
}

type DevImageFlags struct {
	// ExampleFlag string
	Branch string
}

type RuntimeImageFlags struct {
	Branch      string
	MetacallCli bool
}

func (i *DepsImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Branch, "branch", "b", "develop", "eg. of script specific flags")
}

func (i *DevImageFlags) Set(cmd *cobra.Command) {
	// cmd.Flags().StringVarP(&i.ExampleFlag, "exampleb", "e", "", "eg. of script specific flags")
	cmd.Flags().StringVarP(&i.Branch, "branch", "b", "develop", "eg. of script specific flags")
}

func (i *RuntimeImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Branch, "branch", "b", "develop", "eg. of script specific flags")
	cmd.Flags().BoolVar(&i.MetacallCli, "cli", false, "set to also get metacall cli in the runtime image")
}
