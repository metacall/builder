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
			if o.RuntimeImageFlags.MetacallCli {
				args = append(args, "node")
			}
			base := cmd.Context().Value(baseKey{}).(llb.State)

			devBaseEmpty := staging.DevBase(base, branch, []string{})

			devBase := staging.MergeStates(devBaseEmpty)

			runtimeBase := staging.RuntimeBase(base, branch, args)
			finalImage := staging.MergeStates(runtimeBase)

			if o.RuntimeImageFlags.MetacallCli {
				finalImage = staging.AddCli(devBase, finalImage)
			}

			runtime, err := o.Run(finalImage)
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
