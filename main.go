package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/moby/buildkit/client/llb"
	"log"
	"os"
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

	// Random instruction for Buildkit in order to test
	goAlpine := llb.Image("docker.io/library/golang:1.17-alpine")

	dt, err := goAlpine.Marshal(context.TODO(), llb.LinuxAmd64)

	if err != nil {
		log.Fatal(err)
	}

	llb.WriteTo(dt, os.Stdout)
}
