package main

import (
	"os"
	"yoru/globals"
	"yoru/info"
	mkpkg "yoru/make"
	sqlpkg "yoru/sql"
)

var Methods map[string]func([]string)

func init() {
	Methods = map[string]func([]string){
		"make": mkpkg.Main,
		"sql":  sqlpkg.Main,
		"info": info.Main,
	}
}

func main() {
	if len(os.Args) < 2 {
		globals.Handle(false)
		return
	}
	
	command := os.Args[1]
	m, ok := Methods[command]
	globals.Handle(ok)
	m(os.Args[2:])
}
