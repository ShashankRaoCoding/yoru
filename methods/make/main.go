package make

import (
	"yoru/methods/make/dir" 
	"yoru/methods/make/file" 
	"fmt" 
)

var Methods = make(map[string]func([]string)error)

func init() {
	Methods["file"] = file.Main
	Methods["dir"] = dir.Main 
}

func Main(args []string) error {
	m, ok := Methods[args[0]]
	if ok == false {
		return fmt.Errorf("The method %s does not exist", args[0])
	}

	return m(args[1:]) 
}












































