package env

import (
	"github.com/moby/buildkit/client/llb"
)

type Env struct {
	state llb.State
}

func New(base llb.State) Env {
	return Env{state: base}
}

func (e Env) Base() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install git")).
		Root()

	return e
}

func (e Env) MetaCallClone(branch string) Env {
	e.state = e.state.Run(llb.Shlexf("git clone --depth 1 --single-branch --branch=%v https://github.com/metacall/core.git", branch)).
		Run(llb.Shlex("mkdir core/build")).Root()

	return e

}

func (e Env) MetacallEnvBase(arg string) Env {

	e.state = e.state.
		Run(llb.Shlexf("bash core/tools/metacall-environment.sh %v", arg)).
		Root()

	return e
}

func (e Env) MetaCallConfigure(arg string) Env {
	e.state = e.state.File(llb.Mkdir("core/build", 0777)).
		Run(llb.Shlexf("bash core/tools/metacall-configure.sh %v", arg)).
		Root()

	return e
}

func (e Env) MetaCallBuild() Env {
	e.state = e.state.Run(llb.Shlex("bash core/tools/metacall-build.sh")).
		Root()

	return e
}

func (e Env) Root() llb.State {
	return e.state
}
