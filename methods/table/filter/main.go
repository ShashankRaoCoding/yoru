package methods

import (
	"fmt"
	"yoru/methods/table/structs"
)

var Methods = make(map[string]func([]string) error)

func init() {
	Methods["rows"] = filterRows
	Methods["cols"] = filterCols

}

func Main(args []string) error {
	m, ok := Methods[args[0]]
	if ok == false {
		return fmt.Errorf("The method %s does not exist", args[0])
	}

	return m(args[1:])
}

func filterRows([]string) error {
	t, err := structs.New(os.Stdin)

}
