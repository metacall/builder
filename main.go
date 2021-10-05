package main

import (
	"flag"
	"fmt"
	"reflect"
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

var LanguageKeys = reflect.ValueOf(LanguageMap).MapKeys()

func (lang LanguageType) String() string {
	return LanguageKeys[lang].String()
}

type BuildOptions struct {
	version   string
	languages []LanguageType
}

func ParseLanguages(args []string) ([]LanguageType, error) {
	languages := make([]LanguageType, len(args))

	for i, lang := range args {
		langId, found := LanguageMap[lang]

		if found != true {
			return nil, fmt.Errorf("Invalid language, the key '%s' does not exist", lang)
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
		panic(err)
	}

	// Debug
	fmt.Println(opt.version)
	fmt.Println(opt.languages)
}
