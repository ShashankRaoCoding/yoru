package utils

import (
	"fmt"
	"os"
)

// GetBinaryPath returns the absolute path to the currently running executable.
func GetBinaryPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("getting executable path: %w", err)
	}
	return path, nil
}

func Error(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error: %s\n", err)
		os.Exit(1)
	}
}
