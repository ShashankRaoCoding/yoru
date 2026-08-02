package info

import (
	"fmt"
)

func Main(args []string) {
	fmt.Println("yoru - A personal library of binaries")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  make  - Create new files with empty lines")
	fmt.Println("  sql   - Query data with format conversion")
	fmt.Println("  sixel - Convert and display images in SIXEL format")
	fmt.Println("  print - Render images as SIXEL or tabular data as a pretty table")
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
	fmt.Println("  image            Render stdin image data to SIXEL")
	fmt.Println("  table            Render stdin CSV/TSV data to a pretty table")
	fmt.Println("\nSIXEL Examples:")
	fmt.Println("  yoru sixel compatible")
	fmt.Println("\nPrint Examples:")
	fmt.Println("  cat image.png | yoru print image -i png")
	fmt.Println("  cat image.jpg | yoru print image -i jpeg 100 50")
	fmt.Println("  cat data.csv | yoru print table -i csv")
	fmt.Println("  cat data.tsv | yoru print table -i tsv")
}
