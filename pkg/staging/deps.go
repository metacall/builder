package staging

import (
	"fmt"
	"github.com/metacall/builder/pkg/env"
	"github.com/moby/buildkit/client/llb"
	"strings"
)

var languageMap = map[string]string{
	// "cache": "cache",
	"base": "base",
	"py":   "python",
	"rb":   "ruby",
	// "netcore":   "netcore",
	// "netcore2":  "netcore2",
	// "netcore5":  "netcore5",
	// "netcore7":  "netcore7",
	// "rapidjson": "rapidjson",
	// "funchook":  "funchook",
	// "v8":        "v8",
	// "v8rep54":   "v8rep54",
	// "v8rep57":   "v8rep57",
	// "v8rep58":   "v8rep58",
	// "v8rep52":   "v8rep52",
	// "v8rep51":   "v8rep51",
	"nodejs":     "nodejs",
	"typescript": "typescript",
	// "file":       "file",
	// "rpc":        "rpc",
	"wasm":  "wasm",
	"java":  "java",
	"c":     "c",
	"cobol": "cobol",
	"go":    "go",
	"rust":  "rust",
	// "swig":  "swig",
	// "pack":     "pack",
	// "coverage": "coverage",
	// "clangformat": "clangformat",
	// "backtrace"	: "backtrace",
	// "sandbox"	: "sandbox",
	// "scripts":    "scripts",
	// "examples":   "examples",
	// "tests":      "tests",
	// "benchmarks": "benchmarks",
	// "ports":      "ports",
}

func DepsBase(base llb.State, branch string, args []string) llb.State {
	err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		//handle error
		return base
	}

	cmdArgs := strings.Join(args, " ")

	return env.New(base).
		Base().
		MetaCallClone(branch).
		MetacallEnvBase(cmdArgs).
		Root()
}

func DevBase(base llb.State, branch string, args []string) llb.State {
	err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		// TODO : handle error
		return base
	}

	cmdArgs := strings.Join(args, " ")

	return env.New(base).
		Base().
		MetaCallClone(branch).
		MetacallEnvBase(cmdArgs).
		MetaCallConfigure(cmdArgs).
		MetaCallBuild().
		Root()
}

func RuntimeBase(base llb.State, branch string, args []string) llb.State {
	err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		// TODO : handle error
		return base
	}

	cmdArgs := strings.Join(args, " ")

	return env.New(base).
		Base().
		MetaCallClone(branch).
		MetacallRuntime(cmdArgs).
		Root()
}

func AddCli(src llb.State, dst llb.State) llb.State {
	return dst.With(copyFrom(src, "/usr/local/bin/metacallcli*", "/usr/local/bin/metacall"))
}

func CopyFromBuilder(src llb.State, dst llb.State) llb.State {

	// libraries
	libPaths := []string{"/usr/local/lib/*.so", "/usr/local/lib/*.so*", "/usr/local/lib/*.dll", "/usr/local/lib/*.js", "/usr/local/lib/*.ts", "/usr/local/lib/*.node"}
	libDst := "/usr/local/lib/"

	// plugins
	pluginsPath := []string{"/usr/local/lib/plugins"}
	pluginsDst := "/usr/local/lib/plugins"

	// node dependencies (and port)
	ndpPath := []string{"/usr/local/lib/node_modules/"}
	ndpDst := "/usr/local/lib/node_modules/"

	// python dependencies
	pydPath := []string{"/usr/local/lib/python3.11/dist-packages/metacall/"}
	pydDst := "/usr/local/lib/python3.11/dist-packages/metacall/"

	// headers
	hdPath := []string{"/usr/local/include/metacall"}
	hdDst := "/usr/local/include/metacall"

	// configurations
	configPath := []string{"/usr/local/share/metacall/configurations/*"}
	configDst := "/usr/local/share/metacall/configurations/"

	return dst.With(
		copyMultiple(src, libPaths, libDst),
		copyMultiple(src, pluginsPath, pluginsDst),
		copyMultiple(src, ndpPath, ndpDst),
		copyMultiple(src, pydPath, pydDst),
		copyMultiple(src, hdPath, hdDst),
		copyMultiple(src, configPath, configDst),
	)
}
