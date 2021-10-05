package main

import (
	"fmt"
	"flag"
)

type LanguageType string

const (
	Python LanguageType = "py"
	Node = "node"
	TypeScript = "ts"
	Ruby = "rb"
	CSharp = "cs"
	Java = "java"
	WebAssembly = "wasm"
	Cobol = "cob"
	File = "file"
	Mock = "mock"
	RPC = "rpc"
)

type BuildOptions struct {
	version string
	languages []LanguageType
}

func main() {
	var opt BuildOptions
	flag.StringVar(&opt.version, "version", "v0.5.6", "MetaCall version to build with")
	flag.Parse()
	opt.languages = flag.Args() // TODO

	// Debug
	fmt.Println(opt.version)
	fmt.Println(opt.languages)
}
