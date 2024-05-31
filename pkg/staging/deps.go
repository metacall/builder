package staging

import (
	"fmt"
	"strings"

	"github.com/metacall/builder/pkg/env"
	"github.com/moby/buildkit/client/llb"
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
	// TODO : validate args having some bugs

	// _ := validateArgs(args)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// TODO : handle error
	// 	return base
	// }

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
		// TODO : handle error in a better way
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

func RemoveBuild(state llb.State) llb.State {
	return state.File(llb.Rm("build"))
}
