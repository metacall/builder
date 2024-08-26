package builder

import (
	"github.com/spf13/cobra"
)

type RuntimeImageFlags struct {
	MetacallCli bool
	Startup bool
}

func (i *RuntimeImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&i.MetacallCli, "cli", false, "set to also get metacall cli in the runtime image")
}

func (i *RuntimeImageFlags) SetStartup(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&i.Startup, "startup", false, "startup flag to be used for building image with all languages")
}