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
	e.state = e.state.
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
	e.state = e.state.
		AddEnv(
			"LOADER_LIBRARY_PATH",
			"/usr/local/metacall/build",
		).AddEnv(
			"LOADER_SCRIPT_PATH",
			"/usr/local/metacall/build/scripts",
		).AddEnv(
			"CONFIGURATION_PATH",
			"/usr/local/metacall/build/configurations/global.json",
		).AddEnv(
			"SERIAL_LIBRARY_PATH",
			"/usr/local/metacall/build",
		).AddEnv(
			"DETOUR_LIBRARY_PATH",
			"/usr/local/metacall/build",
		).AddEnv(
			"PORT_LIBRARY_PATH",
			"/usr/local/metacall/build",
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
	e.state = e.state.
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
		Run(llb.Shlex("apt-get -y --no-install-recommends install git ca-certificates")).
		Run(llb.Shlex("apt install npm")). // Maybe a better option would be to make a check in metacall environment installation and make langauge specific some installs
		Root()

	return e
}

func (e Env) MetaCallClone(branch string) Env {
	a := e.state
	a = a.Run(llb.Shlexf("git clone --depth 1 --single-branch --branch=%v https://github.com/metacall/core.git", branch)).
	Root()

	e.state = copy(a, "/core/", e.state, "/usr/local/metacall")

	return e
}

func (e Env) MetacallEnvBase(arg string) Env {
	e.state = e.state.
		Run(llb.Shlexf("bash /usr/local/metacall/tools/metacall-environment.sh relwithdebinfo base backtrace swig %v", arg)).Root()
	return e
}

func (e Env) MetaCallConfigure(arg string) Env {
	e.state = e.state.File(llb.Mkdir("build", 0777)).Dir("/usr/local/metacall/build").
		Run(llb.Shlexf("bash /usr/local/metacall/tools/metacall-configure.sh relwithdebinfo tests scripts ports install %v", arg)).
		Root()

	return e
}

func (e Env) MetaCallBuild(arg string) Env {
	e.state = e.state.Run(llb.Shlexf("bash /usr/local/metacall/tools/metacall-build.sh relwithdebinfo tests scripts ports install %v",arg)).
		Root()

	return e
}

func (e Env) MetacallRuntime(arg string) Env {
	e.state = e.state.Run(llb.Shlexf("bash /usr/local/metacall/tools/metacall-runtime.sh base backtrace ports clean %v", arg)).
		Root()

	return e
}

func (e Env) Root() llb.State {
	return e.state
}

func copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	return dest.File(llb.Copy(src, srcPath, destPath, &llb.CopyInfo{
		AllowWildcard:  true,
		AttemptUnpack:  true,
		CreateDestPath: true,
	}))
}
