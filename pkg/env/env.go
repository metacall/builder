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

func (e Env) DepsEnv() Env {
	e.state = e.state.Dir("/usr/local/metacall").
		AddEnv(
			"DEBIAN_FRONTEND",
			"noninteractive",
		).AddEnv(
		"LTTNG_UST_REGISTER_TIMEOUT",
		"0",
	).AddEnv(
		"NUGET_XMLDOC_MODE",
		"skip",
	).AddEnv(
		"DOTNET_CLI_TELEMETRY_OPTOUT",
		"true",
	)
	return e
}

func (e Env) Base() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install git")).
		Root()

	return e
}

func (e Env) MetaCallClone(branch string) Env {
	e.state = e.state.Run(llb.Shlexf("git -c http.sslVerify=false clone --depth 1 --single-branch --branch=%v https://github.com/metacall/core.git", branch)).
		Root().Dir("/usr/local/metacall/core")
	return e

}

func (e Env) MetacallEnvBase(arg string) Env {
	e.state = e.state.
		Run(llb.Shlexf("bash tools/metacall-environment.sh relwithdebinfo base backtrace %v", arg)).
		// Dir("/").Run(llb.Shlex("rm -rf /usr/local/metacall")).
		Root()
	return e
}

func (e Env) MetaCallConfigure(arg string) Env {
	e.state = e.state.File(llb.Mkdir("build", 0777)).
		Run(llb.Shlexf("bash tools/metacall-configure.sh base backtrace tests scripts ports install %v", arg)).
		Root()

	return e
}

func (e Env) MetaCallBuild() Env {
	e.state = e.state.Run(llb.Shlex("bash tools/metacall-build.sh")).
		Root()

	return e
}

func (e Env) MetacallRuntime(arg string) Env {
	e.state = e.state.Run(llb.Shlexf("bash tools/metacall-runtime.sh backtrace ports %v", arg)).
		Root()

	return e
}

func (e Env) Root() llb.State {
	return e.state
}
