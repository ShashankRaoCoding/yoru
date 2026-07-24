package info

import (
	"fmt"
)

func Main(args []string) {
	fmt.Println("yoru - A personal library of binaries")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  make  - Create new files with empty lines")
	fmt.Println("  sql   - Query data with format conversion")
	fmt.Println("  info  - Display this help message")
	fmt.Println("\nSQL Command Usage:")
	fmt.Println("  yoru sql [options] <query>")
	fmt.Println("\nSQL Options:")
	fmt.Println("  -i format   Input format: csv, tsv, or sqldb (default: csv)")
	fmt.Println("  -o format   Output format: csv or tsv (default: csv)")
	fmt.Println("\nSQL Examples:")
	fmt.Println("  cat data.csv | yoru sql -i csv -o csv 'SELECT * FROM table'")
	fmt.Println("  cat data.tsv | yoru sql -i tsv -o tsv 'SELECT * FROM table'")
	fmt.Println("  printf 'input.db' | yoru sql -i sqldb -o csv 'SELECT * FROM table'")
}
