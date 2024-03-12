package builder

import (
	"github.com/spf13/cobra"
)

type DepsImageFlags struct {
	Branch string
}

type DevImageFlags struct {
	ExampleFlag string
}

type RuntimeImageFlags struct {
	ExampleFlag string
}

func (i *DepsImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Branch, "branch", "b", "develop", "eg. of script specific flags")
}

func (i *DevImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.ExampleFlag, "exampleb", "exb", "", "eg. of script specific flags")
}

func (i *RuntimeImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.ExampleFlag, "examplec", "exc", "", "eg. of script specific flags")
}
