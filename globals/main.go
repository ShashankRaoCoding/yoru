package globals

import (
	"fmt"
	"os"
	"path/filepath"
)

var TEMP string

func init() {
	// Try to use system temp directory
	temp := os.TempDir()
	TEMP = filepath.Join(temp, "yoru")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(TEMP, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error creating temp directory, using ./temp instead\n")
		TEMP = "./temp"
		os.MkdirAll(TEMP, 0755)
	}
}
