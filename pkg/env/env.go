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

func (e Env) DevEnv() Env {
	e.state = e.state.Dir("/usr/local/metacall/core").
		AddEnv(
			"LOADER_LIBRARY_PATH",
			"/usr/local/metacall/core/build",
		).AddEnv(
			"LOADER_SCRIPT_PATH",
			"/usr/local/metacall/core/build/scripts",
		).AddEnv(
			"CONFIGURATION_PATH",
			"/usr/local/metacall/core/build/configurations/global.json",
		).AddEnv(
			"SERIAL_LIBRARY_PATH",
			"/usr/local/metacall/core/build",
		).AddEnv(
			"DETOUR_LIBRARY_PATH",
			"/usr/local/metacall/core/build",
		).AddEnv(
			"PORT_LIBRARY_PATH",
			"/usr/local/metacall/core/build",
		).AddEnv(
			"NODE_PATH",
			"/usr/lib/node_modules",
		).AddEnv(
			"DEBIAN_FRONTEND",
			"noninteractive",
		).AddEnv(
			"DOTNET_CLI_TELEMETRY_OPTOUT",
			"true",
		).File(llb.Mkdir("build", 0755))
	return e
}

func (e Env) RuntimeEnv() Env {
	e.state = e.state.Dir("/usr/local/metacall").
	AddEnv(
		"DEBIAN_FRONTEND",
		"noninteractive",
	).AddEnv(
		"DOTNET_CLI_TELEMETRY_OPTOUT",
		"true",
	).AddEnv(
		"LOADER_LIBRARY_PATH",
		"/usr/local/li",
	).AddEnv(
		"LOADER_SCRIPT_PATH",
		"/usr/local/scripts",
	).AddEnv(
		"CONFIGURATION_PATH",
		"/usr/local/share/metacall/configurations/global.json",
	).AddEnv(
		"SERIAL_LIBRARY_PATH",
		"/usr/local/lib",
	).AddEnv(
		"DETOUR_LIBRARY_PATH",
		"/usr/local/lib",
	).AddEnv(
		"PORT_LIBRARY_PATH",
		"/usr/local/lib",
	).AddEnv(
		"NODE_PATH",
		"/usr/local/lib/node_modules",
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
		Run(llb.Shlexf("bash /usr/local/metacall/core/tools/metacall-environment.sh relwithdebinfo base backtrace %v", arg)).Root()
	return e
}

func (e Env) MetaCallConfigure(arg string) Env {
	e.state = e.state.File(llb.Mkdir("build", 0777)).Dir("/usr/local/metacall/core/build").
		Run(llb.Shlexf("bash /usr/local/metacall/core/tools/metacall-configure.sh relwithdebinfo tests scripts ports install %v", arg)).
		Root()

	return e
}

func (e Env) MetaCallBuild(arg string) Env {
	e.state = e.state.Run(llb.Shlexf("bash /usr/local/metacall/core/tools/metacall-build.sh relwithdebinfo %v",arg)).
		Root()

	return e
}

func (e Env) MetacallRuntime(arg string) Env {
	e.state = e.state.Run(llb.Shlexf("bash /usr/local/metacall/core/tools/metacall-runtime.sh base backtrace ports clean %v", arg)).
		Root()

	return e
}

func (e Env) Root() llb.State {
	return e.state
}
