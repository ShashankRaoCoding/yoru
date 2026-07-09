package table

import (
	"bufio"
	"os"
)

var Methods = make(map[string]func([]string) error)

func init() {
	Methods[""]

}

func Main(args []string) error {

}
