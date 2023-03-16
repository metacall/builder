package builder

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "builder",
		Short:         "builder is a tool for building MetaCall images",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(NewDoctorCmd())

	return cmd
}
