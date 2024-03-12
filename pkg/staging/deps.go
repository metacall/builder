package staging

import (
	"github.com/metacall/builder/pkg/env"
	"github.com/moby/buildkit/client/llb"
)

type Deps struct{}

func (deps *Deps) Base(base llb.State, branch string) llb.State {
	return env.New(base).
		Base().
		MetaCallClone(branch).
		MetacallEnvBase().
		Root()
}
