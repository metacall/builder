package builder

import (
	"context"

	"github.com/metacall/builder/pkg/staging"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

type DepOptions struct {
	// DepsImageFlags DepsImageFlags
}

func NewDepsOptions() *DepOptions {
	return &DepOptions{}
}

func NewDepsCmd(o *DepOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deps",
		Short: "Build development dependencies base image for MetaCall",
		RunE: func(cmd *cobra.Command, args []string) error {
			base := cmd.Context().Value(baseKey{}).(llb.State)
			depsBase := staging.DepsBase(base, branch, args)
			depsBase, err := o.Run(depsBase)
			if err != nil {
				return err
			}

			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, depsBase))
			return nil

		},
		Example: `builder deps -b develop nodejs typescript go rust wasm java c cobol`,
	}
	// o.DepsImageFlags.Set(branch)

	return cmd
}

func (do *DepOptions) Run(depsBase llb.State) (llb.State, error) {
	// return depsBase.Dir("/").Run(llb.Shlex("rm -rf /usr/local/metacall")).Root(), nil
	return depsBase, nil
}
