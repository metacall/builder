package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/moby/buildkit/client/llb"
)

type LanguageType int

const (
	Python LanguageType = iota
	Node
	TypeScript
	Ruby
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
	"cs":   CSharp,
	"java": Java,
	"wasm": WebAssembly,
	"cob":  Cobol,
	"file": File,
	"mock": Mock,
	"rpc":  RPC,
}

var buildFuncMap = map[LanguageType]func(llb.State) llb.State{
	Python: buildPy,
	Node:   buildNode,
	//	TypeScript: buildTS,
	Ruby:        buildRuby,
	CSharp:      buildCSharp,
	Java:        buildJava,
	WebAssembly: buildWasm,
	Cobol:       buildCobol,
	//	File:        buildFile,
	RPC: buildRPC,
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

func buildNode(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update && apt-get -y --no-install-recommends install python3 g++ make nodejs curl")).Root()
}

func buildPy(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update && apt-get -y --no-install-recommends install python3 python3-dev python3-pip")).Root()
}

func buildRuby(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update && apt-get -y --no-install-recommends install ruby2.7 ruby2.7-dev")).Root()
}

func buildRPC(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update && apt-get -y --no-install-recommends install libcurl4-openssl-dev")).Root()
}

func buildJava(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update && apt-get install default-jdk")).Root()
}

func buildCobol(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex("apt-get update && apt-get install open-cobol")).Root()
}

func buildCSharp(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex(`apt-get update && apt-get install -y apt-transport-https && \
		apt-get update && \
		apt-get install -y dotnet-sdk-5.0`)).Root()
}

func buildWasm(baseImg llb.State) llb.State {
	return baseImg.
		Run(llb.Shlex(`apt-get update && apt-get install curl ca-certificates \
		&& apt-get install xz-utils \
		&& curl https://wasmtime.dev/install.sh -sSf | bash\
		&& apt-get remove curl ca-certificates \
		&& apt-get remove xz-utils \
		88 apt-get autoremove --purge`)).Root()
}

func buildDeps(langs []LanguageType) {

	// Pulls Debian BaseImage from registry
	baseImg := llb.Image("docker.io/library/debian:bullseye-slim")

	dt, err := baseImg.Marshal(context.TODO(), llb.LinuxAmd64)

	for _, v := range langs {
		baseImg = buildFuncMap[v](baseImg)
	}
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

	// Debug
	fmt.Println(opt.version)
	fmt.Println(opt.languages)

	buildDeps(opt.languages)

}
