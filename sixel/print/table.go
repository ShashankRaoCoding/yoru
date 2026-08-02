package print

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func runTable(args []string) error {
	fs := flag.NewFlagSet("table", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	inputFormat := fs.String("i", "csv", "Input format: csv or tsv")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: yoru sixel print table -i <format>\n")
		fmt.Fprintf(os.Stderr, "\nInput formats:\n")
		fmt.Fprintf(os.Stderr, "  csv, tsv\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  cat data.csv | yoru sixel print table -i csv\n")
		fmt.Fprintf(os.Stderr, "  cat data.tsv | yoru sixel print table -i tsv\n")
	}

	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected extra arguments")
	}

	format := strings.ToLower(strings.TrimSpace(*inputFormat))
	delimiter := ','
	switch format {
	case "csv":
	case "tsv":
		delimiter = '\t'
	default:
		return fmt.Errorf("unsupported input format %q", *inputFormat)
	}

	out, err := renderDelimitedAsTable(os.Stdin, delimiter)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, out)
	return err
}

func renderDelimitedAsTable(r io.Reader, delimiter rune) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("reading input: %w", err)
	}

	reader := csv.NewReader(strings.NewReader(string(input)))
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("reading delimited input: %w", err)
	}
	if len(records) == 0 {
		return "", nil
	}

	maxCols := 0
	for _, row := range records {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	widths := make([]int, maxCols)
	rows := make([][]string, len(records))
	for i, row := range records {
		normalized := make([]string, maxCols)
		copy(normalized, row)
		rows[i] = normalized
		for colIdx, cell := range normalized {
			if len(cell) > widths[colIdx] {
				widths[colIdx] = len(cell)
			}
		}
	}

	border := func(left, mid, right string) string {
		var b strings.Builder
		b.WriteString(left)
		for i, width := range widths {
			b.WriteString(strings.Repeat("-", width+2))
			if i < len(widths)-1 {
				b.WriteString(mid)
			}
		}
		b.WriteString(right)
		return b.String()
	}

	var out strings.Builder
	out.WriteString(border("+", "+", "+"))
	out.WriteString("\n")

	for rowIdx, row := range rows {
		out.WriteString("|")
		for colIdx, cell := range row {
			padding := widths[colIdx] - len(cell)
			out.WriteString(" ")
			out.WriteString(cell)
			out.WriteString(strings.Repeat(" ", padding+1))
			out.WriteString("|")
		}
		out.WriteString("\n")

		if rowIdx == 0 && len(rows) > 1 {
			out.WriteString(border("+", "+", "+"))
			out.WriteString("\n")
		}
	}

	out.WriteString(border("+", "+", "+"))
	return out.String(), nil
}
