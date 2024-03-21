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
			devBase := staging.DevBase(base, o.RuntimeImageFlags.Branch, args)
			runtimeBase := staging.RuntimeBase(base, o.RuntimeImageFlags.Branch, args)

			runtime := staging.CopyFromBuilder(devBase, runtimeBase)
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
		Example: `"builder runtime -b runtime --cli nodejs typescript go rust wasm java c cobol"`,
	}
	o.RuntimeImageFlags.Set(cmd)

	return cmd
}

func (do *RuntimeOptions) Run(runtimeBase llb.State) (llb.State, error) {
	return runtimeBase, nil
}
