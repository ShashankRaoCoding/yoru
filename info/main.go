package info

import (
	"fmt"
)

func Main(args []string) {
	fmt.Println("yoru - A personal library of binaries")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  make  - Create new files with empty lines")
	fmt.Println("  print - Render images in SIXEL format or print tables")
	fmt.Println("  sql   - Query data with format conversion")
	fmt.Println("  sixel - Check SIXEL terminal compatibility")
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
	fmt.Println("\nSIXEL Subcommands:")
	fmt.Println("  compatible       Check if the terminal supports SIXEL graphics")
	fmt.Println("\nPrint Subcommands:")
	fmt.Println("  image            Render images to SIXEL")
	fmt.Println("  table            Render tabular text to a pretty table")
	fmt.Println("\nExamples:")
	fmt.Println("  yoru sixel compatible")
	fmt.Println("  cat image.png | yoru print image -i png")
	fmt.Println("  cat image.jpg | yoru print image -i jpeg 100 50")
	fmt.Println("  cat data.csv | yoru print table -i csv")
}
