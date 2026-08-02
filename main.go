package main

import (
	"fmt"
	"os"
	"yoru/info"
	printpkg "yoru/print"
	"yoru/utils"
	mkpkg "yoru/make"
	sixelpkg "yoru/sixel"
	sqlpkg "yoru/sql"
)

var Methods map[string]func([]string)

func init() {
	Methods = map[string]func([]string){
		"make":  mkpkg.Main,
		"print": printpkg.Main,
		"sql":   sqlpkg.Main,
		"sixel": sixelpkg.Main,
		"info":  info.Main,
	}
}

func main() {
	if len(os.Args) < 2 {
		utils.Error(fmt.Errorf("no command provided"))
		return
	}
	
	command := os.Args[1]
	m, ok := Methods[command]
	if !ok {
		utils.Error(fmt.Errorf("method not found: %s", command))
		return
	}
	m(os.Args[2:])
}
