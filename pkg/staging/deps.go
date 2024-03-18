package staging

import (
	"errors"
	"fmt"
	"github.com/metacall/builder/pkg/env"
	"github.com/moby/buildkit/client/llb"
	"strings"
)

type Deps struct{}

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

func ValidateArgs(args []string) error {
	for _, arg := range args {
		if _, ok := languageMap[arg]; !ok {
			return errors.New("Invalid language")
		}
	}
	return nil
}

func (deps *Deps) DevBase(base llb.State, branch string, args []string) llb.State {
	err := ValidateArgs(args)
	if err != nil {
		fmt.Println(err)
		//handle error
		return base
	}

	cmdArgs := strings.Join(args, " ")

	fmt.Println(cmdArgs)

	return env.New(base).
		Base().
		MetaCallClone(branch).
		MetacallEnvBase(cmdArgs).
		MetaCallConfigure(cmdArgs).
		MetaCallBuild().
		Root()
}

func (deps *Deps) DepsBase(base llb.State, branch string, args []string) llb.State {

	err := ValidateArgs(args)
	if err != nil {
		fmt.Println(err)
		//handle error
		return base
	}

	cmdArgs := strings.Join(args, " ")

	fmt.Println(cmdArgs)

	return env.New(base).
		Base().
		MetaCallClone(branch).
		MetacallEnvBase(cmdArgs).
		Root()

}
