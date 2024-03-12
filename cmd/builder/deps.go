package builder

import (
	"context"
	"github.com/metacall/builder/pkg/staging"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

type DepOptions struct {
	//TODO : add ui using cli ui's. (example is "github.com/cppforlife/go-cli-ui/ui") or others
	DepsImageFlags DepsImageFlags
}

// include the ui cli part here
func NewDepsOptions() *DepOptions {
	return &DepOptions{}
}

func NewDepsCmd(o *DepOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deps",
		Short: "Build development dependencies base image for MetaCall",
		RunE: func(cmd *cobra.Command, args []string) error {
			base := cmd.Context().Value(baseKey{}).(llb.State)
			myDeps := staging.Deps{}
			depsBase := myDeps.Base(base, o.DepsImageFlags.Branch)
			depsBase, err := o.Run(depsBase)
			if err != nil {
				return err
			}
			// set final state
			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, depsBase))
			return nil

		},
		Example: `blablalablabol`,
	}
	o.DepsImageFlags.Set(cmd)

	return cmd
}

func (do *DepOptions) Run(depsBase llb.State) (llb.State, error) {
	return depsBase, nil
}
