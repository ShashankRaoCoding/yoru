package main

import (
	"fmt"
	"os"
	"yoru/info"
	"yoru/utils"
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
