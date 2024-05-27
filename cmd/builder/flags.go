package builder

import (
	"github.com/spf13/cobra"
)

type RuntimeImageFlags struct {
	MetacallCli bool
}

func (i *RuntimeImageFlags) Set(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&i.MetacallCli, "cli", false, "set to also get metacall cli in the runtime image")
}
