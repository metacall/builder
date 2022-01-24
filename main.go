package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/moby/buildkit/client/llb"
	_ "github.com/moby/buildkit/util/progress"
)

type LanguageType int

const (
	Python LanguageType = iota
	Node
	TypeScript
	Ruby
	C
	CSharp
	Java
	WebAssembly
	Cobol
	File
	Mock
	RPC
)

var LanguageMap = map[string]LanguageType{
	"py":   Python,
	"node": Node,
	"ts":   TypeScript,
	"rb":   Ruby,
	"c":    C,
	"cs":   CSharp,
	"java": Java,
	"wasm": WebAssembly,
	"cob":  Cobol,
	"file": File,
	"mock": Mock,
	"rpc":  RPC,
}

var buildFuncMap = map[LanguageType]func(llb.State) llb.State{
	Python: buildPyEnv,
	Node:   buildNodeEnv,
	//	TypeScript: buildTS,
	Ruby:        buildRubyEnv,
	CSharp:      buildCSharpEnv,
	Java:        buildJavaEnv,
	WebAssembly: buildWasmEnv,
	C:           buildCEnv,
	//	Cobol:       buildCobol,
	//	File:        buildFile,
	RPC: buildRPCEnv,
}

var LanguageKeys = (func() []string {
	keys := make([]string, len(LanguageMap))
	for k, v := range LanguageMap {
		keys[v] = k
	}
	return keys
})()

func (lang LanguageType) String() string {
	return LanguageKeys[lang]
}

type BuildOptions struct {
	version   string
	languages []LanguageType
}

func ParseLanguages(args []string) ([]LanguageType, error) {
	languages := make([]LanguageType, len(args))

	for i, lang := range args {
		langId, found := LanguageMap[lang]

		if !found {
			return nil, fmt.Errorf("invalid language, the key '%s' does not exist", lang)
		}

		languages[i] = langId
	}

	return languages, nil
}

func buildNodeEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install python3 g++ make nodejs curl")).Root()
}

func buildPyEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install python3 python3-dev python3-pip")).Root()
}

func buildRubyEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install ruby2.7 ruby2.7-dev")).Root()
}

func buildRPCEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install libcurl4-openssl-dev")).Root()
}

func buildJavaEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install default-jdk")).Root()
}

func buildCSharpEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install wget")).
		Run(llb.Shlex("wget https://packages.microsoft.com/config/debian/11/packages-microsoft-prod.deb -O packages-microsoft-prod.deb")).
		Run(llb.Shlex("dpkg -i packages-microsoft-prod.deb")).
		Run(llb.Shlex("rm packages-microsoft-prod.deb")).
		Run(llb.Shlex("apt-get install -y apt-transport-https")).
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y dotnet-sdk-5.0")).
		Run(llb.Shlex("apt-get -y remove wget")).
		Run(llb.Shlex("apt-get -y autoremove --purge")).Root()

}

func buildWasmEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install xz-utils wget")).
		Run(llb.Shlex("wget https://wasmtime.dev/install.sh")).
		Run(llb.Shlex("bash install.sh")).
		Run(llb.Shlex("rm install.sh")).
		Run(llb.Shlex("apt-get -y remove wget xz-utils")).
		Run(llb.Shlex("apt-get -y autoremove --purge")).Root()
}

func buildCEnv(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install cmake build-essential")).Root()
}

func buildDeps(langs []LanguageType) {

	// Pulls Debian BaseImage from registry
	baseImg := llb.Image("docker.io/library/debian:bullseye-slim")

	for _, v := range langs {
		baseImg = buildFuncMap[v](baseImg)
	}
	dt, err := baseImg.Marshal(context.TODO(), llb.LinuxAmd64)

	if err != nil {
		log.Fatal(err)
	}

	llb.WriteTo(dt, os.Stdout)
}

func main() {
	var opt BuildOptions
	var err error
	flag.StringVar(&opt.version, "version", "v0.5.6", "MetaCall version to build with")
	flag.Parse()
	opt.languages, err = ParseLanguages(flag.Args())

	if err != nil {
		log.Fatal(err)
	}

	buildDeps(opt.languages)

}
