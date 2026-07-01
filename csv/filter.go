package csv

import (
	"fmt"
	"os"
	"os/exec"
)

func Init() {
	Methods["filter"] = filter
}

func filter(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("yoru csv filter: expected a method (rows|cols)")
	}

	method, rest := args[0], args[1:]

	m, ok := map[string]func([]string) error{
		"rows": filterRows,
		"cols": filterCols,
	}[method]

	if !ok {
		return fmt.Errorf("Error: Method %s does not exist for yoru csv filter", method)
	}

	return m(rest)
}

// format: shanks csv filter rows columnName condition
// e.g.:   shanks csv filter rows age ">30"
func filterRows(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("yoru csv filter rows: expected <columnName> <condition>")
	}

	columnName := args[0]
	condition := args[1]

	cmd := exec.Command(
		"python",
		"-c",
		fmt.Sprintf(rowsScript, columnName, condition),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yoru csv filter: There was an error: %s", err)
	}

	return nil
}

// format: shanks csv filter cols col1,col2,col3
func filterCols(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("yoru csv filter cols: expected a comma-separated list of column names")
	}

	cols := strings.Split(args[0], ",")
	for i, c := range cols {
		cols[i] = fmt.Sprintf("'%s'", strings.TrimSpace(c))
	}
	colList := "[" + strings.Join(cols, ", ") + "]"

	cmd := exec.Command(
		"python",
		"-c",
		fmt.Sprintf(colsScript, colList),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yoru csv filter: There was an error: %s", err)
	}

	return nil
}

var rowsScript = `
import sys
import pandas

data = pandas.read_csv(sys.stdin)
data = data[data['%s']%s]
data.to_csv(sys.stdout, index=False)
`

var colsScript = `
import sys
import pandas

data = pandas.read_csv(sys.stdin)
data = data[%s]
data.to_csv(sys.stdout, index=False)
`