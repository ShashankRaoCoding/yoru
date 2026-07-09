package table

import (
	"bufio"
	"os"
)

type Table struct {
	Columns []string
	Dim     [2]int
	Rows    map[string]map[string]*any
	Cols    map[string]map[string]*any
}

func New(data *os.File) *Table {
	var t *Table
	t.Rows = make(map[string]map[string]*any)
	t.Cols = make(map[string]map[string]*any)

}
