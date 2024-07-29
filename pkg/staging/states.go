package staging

import (
	"github.com/metacall/builder/pkg/env"
	"github.com/moby/buildkit/client/llb"
)

var languageMap = map[string]string{
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
	"node": "nodejs",
	"ts":   "typescript",
	// "file":       "file",
	// "rpc":        "rpc",
	"wasm": "wasm",
	"java": "java",
	"c":    "c",
	"cob":  "cobol",
	"go":   "go",
	"rs":   "rust",
	// "backtrace"	: "backtrace",
	// "sandbox"	: "sandbox",
}

func DepsBase(base llb.State, branch string, args []string) map[string]llb.State {
	cmdArgs, err := validateArgs(args)
	if err != nil {
		panic(err)
	}

	m := make(map[string]llb.State, len(cmdArgs))

	envBase := env.New(base).
		DepsEnv().
		Base().
		MetaCallClone(branch)

	for _, arg := range cmdArgs {
		foo := envBase
		depsLang := foo.MetacallEnvBase(arg).Root()
		m[arg] = depsLang
	}

	return m
}

func DevBase(base llb.State, branch string, args []string) map[string]llb.State {

	langMapDev := DepsBase(base, branch, args)

	for lang, langDeps := range langMapDev {

		langDev := env.New(langDeps).
			DevEnv().
			MetaCallConfigure(lang).
			MetaCallBuild(lang).
			Root()

		langMapDev[lang] = langDev
	}

	return langMapDev
}

// func RuntimeBase(base llb.State, branch string, args []string) llb.State {

// 	cmdArgs, err := validateArgs(args)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return env.New(base).
// 		RuntimeEnv().
// 		Base().
// 		MetaCallClone(branch).
// 		MetacallRuntime(cmdArgs).
// 		Root()
// }

func AddCli(src llb.State, dst llb.State) llb.State {
	return dst.With(copyFrom(src, "/usr/local/bin/metacallcli*", "/usr/local/bin/metacallcli"))
}

func RemoveBuild(state llb.State) llb.State {
	return state
	// return state.File(llb.Rm("/usr/local/bin/metacall"))
}

func MergeStates(individualLangStates map[string]llb.State) llb.State {
	states := []llb.State{}
	for _, state := range individualLangStates {
		states = append(states, state)
	}
	return llb.Merge(states)
}
