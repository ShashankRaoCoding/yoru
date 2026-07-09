package structs

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Table struct {
	Columns []string
	Dim     [2]int
	Rows    map[string]map[string]*string
	Cols    map[string]map[string]*string
}

func parseRow(row strings) []string {

	var values []string

	values = strings.Split(text, ",")

	return values
}

func NewTable(source *io.Reader, indexCol string) (*Table, error) {
	var err error
	var t = &Table{}
	var ok bool

	t.Rows = make(map[string]map[string]*string)
	t.Cols = make(map[string]map[string]*string)

	scanner := bufio.NewScanner(source)

	ok = scanner.Scan()
	if ok == false {
		err = scanner.Err()
		if err != nil {
			return t, err
		} else {
			return t, fmt.Errorf("Empty CSV")
		}
	}

	t.Columns = parseRow(scanner.Text())
	for _, c := range t.Columns {
		t.Cols[c] = make(map[string]*string)
	}

	// construct rows
	var values []string
	var index = func(i int) string {
		if indexCol == "" {
			return fmt.Sprintf("%v", i)
		} else {
			return indexCol

		}

	}

	for rowIndex := 0; scanner.Scan(); rowIndex++ {

		values = parseRow(scanner.Text())
		if len(t.Columns) != len(values) {
			return t, fmt.Errorf("Mismatch %v columns and %v values at row %v", len(t.Columns), len(values), index)
		}

		for i, c := range t.Columns {
			row, ok := t.Rows[index(rowIndex)]
			if ok == false {
				row = make(map[string]*string)
				t.Rows[index(rowIndex)] = row
			}

			row[c] = &values[i]
			t.Cols[c][index(rowIndex)] = &values[i]
		}

	}

	if err := scanner.Err(); err != nil {
		return t, err
	}

	t.Dim = [2]int{
		len(t.Rows),
		len(t.Cols),
	}

	return t, nil

}

func (t *Table) String() string {
	var rows []string

	rows = append(rows, strings.Join(t.Columns, ","))
	for _, row := range t.Rows {
		var values []string
		for _, v := range row {
			values = append(values, v)
		}
		row := strings.Join(values, ",")
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func (t *Table) FilterRows(f func(map[string]*string) bool) {
	var ok bool
	for i, row := range t.Rows {
		ok = f(row)
		if ok == false {
			t.DeleteRow(i)
		}
	}
}

func (t *Table) FilterColumns(f func(map[string]*string) bool) {
	var ok bool
	for i, col := range t.Cols {
		ok = f(col)
		if ok == false {
			t.DeleteColumn(i)
		}
	}
}

func (t *Table) DeleteColumn(columnName string) error {
	_, ok := t.Cols[columnName]
	if ok == false {
		return fmt.Errorf("Column with columns name %s does not exist", columnName)
	}

	delete(t.Cols, columnName)
	for _, row := range t.Rows {
		delete(row, columnName)
	}
}

func (t *Table) DeleteRow(rowIndex string) error {
	_, ok := t.Rows[rowIndex]
	if ok == false {
		return fmt.Errorf("Row with index %s does not exist", rowIndex)
	}

	delete(t.Rows, rowIndex)
	for _, c := range t.Cols {
		delete(c, rowIndex)
	}
}
