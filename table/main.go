package main

import (
	mk "make"
	"os"
	"table"
)

var Methods map[string]func([]string)

func init() {
	var Methods = map[string]func([]string){
		"make":  mk.Main,
		"table": table.Main,
		"info":  info.Main,
	}
}

func main() {
	args := os.Args[1:]
	m, ok := Methods[os.Args[0]]
	globals.Handle(ok)
	m(args[1:])
}

func CreateTableFromStdin() {

}
