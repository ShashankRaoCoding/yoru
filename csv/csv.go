package csv

import (
	"os"
	"bufio"
	"os/exec" 
	"yoru/globals" 
)

var Methods = make(map[string]func([]string)error) 

func Init() {
	globals.Methods["csv"] = csv
}

// format: shanks csv filter by columnName where condition
func csv(args []string) error {
	var method, args := args[0], args[1:] 
	
	m, ok := Methods[method] 
	if ok == false {
		return fmt.Errorf("yoru csv: function %s does not exist", method)
	}

	return m(args)

}

var script = `
import sys
import pandas

data = pandas.read_csv(sys.stdin) 
data = data[data[%v]%v]
data.to_csv(sys.stdout, index=False)
`








































