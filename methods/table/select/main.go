package methods

import (
	"fmt"
)

var Methods = make(map[string]func([]string) error)

func init() {
	Methods["filter"] = filter.Main
	Methods["swap"] = replace.Main
	Methods["select"] = select.Main
}

func Main(args []string) error {
	m, ok := Methods[args[0]]
	if ok == false {
		return fmt.Errorf("The method %s does not exist", args[0])
	}

	return m(args[1:])
}
