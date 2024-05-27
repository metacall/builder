package builder

import (
	"context"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

var (
	image  string
	exe   string
	branch string
)

func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:           "builder",
		Short:         "builder is a tool for building MetaCall images",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set base image
			base := llb.Image(image)
			if exe != "" {
				base = base.Run(llb.Shlex(exe)).Root()
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

	cmd.PersistentFlags().StringVarP(&branch, "branch", "b", "develop", "branch to pull metacall from")
	cmd.PersistentFlags().StringVarP(&image, "image", "i", "debian:bullseye-slim", "base image of target image")
	cmd.PersistentFlags().StringVarP(&exe, "exe", "e", "", "exec commands on base image before building (e.g. apt-get update)")

	cmd.AddCommand(NewDoctorCmd())
	cmd.AddCommand(NewDepsCmd(NewDepsOptions()))
	cmd.AddCommand(NewDevCmd(NewDevOptions()))
	cmd.AddCommand(NewRuntimeCmd(NewRuntimeOptions()))

	return cmd
}
