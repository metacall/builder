package builder

import (
	"context"
	"github.com/metacall/builder/pkg/staging"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

const flagBranch = "branch"

func NewDevCmd() *cobra.Command {
	var (
		branch string
	)

	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Build development image for MetaCall",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(NewDevDepsBaseCmd())

	cmd.PersistentFlags().StringVarP(&branch, flagBranch, "b", "develop", "core git branch to use")

	return cmd
}

func NewDevDepsBaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deps-base",
		Short: "Build development dependencies base image for MetaCall",
		RunE: func(cmd *cobra.Command, args []string) error {
			base := cmd.Context().Value(baseKey{}).(llb.State)

			branch, err := cmd.Flags().GetString(flagBranch)
			if err != nil {
				return err
			}
			depsBase := staging.Deps.Base(base, branch)

			// set final state
			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, depsBase))
			return nil
		},
	}

	return cmd
}
