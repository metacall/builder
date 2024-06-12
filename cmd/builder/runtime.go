package builder

import (
	"context"
	
	"github.com/metacall/builder/pkg/staging"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

type RuntimeOptions struct {
	RuntimeImageFlags RuntimeImageFlags
}

func NewRuntimeOptions() *RuntimeOptions {
	return &RuntimeOptions{}
}

func NewRuntimeCmd(o *RuntimeOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runtime",
		Short: "Build runtime image for MetaCall",
		RunE: func(cmd *cobra.Command, args []string) error {
			base := cmd.Context().Value(baseKey{}).(llb.State)
			devBase := staging.RemoveBuild(staging.DevBase(base, branch, []string{}))
			devBaseLang := staging.RemoveBuild(staging.DevBase(base, branch, args))
			runtimeBase := staging.RuntimeBase(base, branch, args)
			diffed := llb.Diff(devBase, devBaseLang)

			runtime := llb.Merge([]llb.State{runtimeBase, diffed})

			if o.RuntimeImageFlags.MetacallCli {
				runtime = staging.AddCli(devBase, runtime)
			}

			runtime, err := o.Run(runtime)
			if err != nil {
				return err
			}

			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, runtime))
			return nil

		},
		Example: `"builder runtime -b develop --cli nodejs typescript go rust wasm java c cobol"`,
	}
	o.RuntimeImageFlags.Set(cmd)

	return cmd
}

func (do *RuntimeOptions) Run(runtimeBase llb.State) (llb.State, error) {
	return runtimeBase, nil
}
