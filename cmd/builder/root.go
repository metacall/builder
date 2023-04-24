package builder

import (
	"context"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
	"os"
)

func NewRootCmd() *cobra.Command {
	var (
		image string
		exec  string
	)

	cmd := &cobra.Command{
		Use:           "builder",
		Short:         "builder is a tool for building MetaCall images",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set base image
			base := llb.Image(image)
			if exec != "" {
				base = base.Run(llb.Shlex(exec)).Root()
			}
			cmd.SetContext(context.WithValue(cmd.Context(), baseKey{}, base))
			// set languages
			cmd.SetContext(context.WithValue(cmd.Context(), languagesKey{}, args))
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			state, ok := cmd.Context().Value(finalKey{}).(llb.State)
			if !ok {
				return nil
			}

			def, err := state.Marshal(cmd.Context(), llb.LinuxAmd64)
			if err != nil {
				return err
			}

			return llb.WriteTo(def, os.Stdout)
		},
	}

	cmd.AddCommand(NewDoctorCmd(), NewDevCmd())

	cmd.PersistentFlags().StringVarP(&image, "image", "i", "debian:bullseye-slim", "base image of target image")
	cmd.PersistentFlags().StringVarP(&exec, "exec", "e", "", "exec commands on base image before building (e.g. apt-get update)")

	return cmd
}
