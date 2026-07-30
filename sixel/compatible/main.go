package compatible

import (
	"fmt"
	"os"
	"strings"
)

// Main checks whether the current terminal supports the SIXEL graphics protocol
// and prints "true" or "false" accordingly.
func Main(args []string) {
	if IsCompatible() {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

// IsCompatible returns true when heuristics suggest the terminal supports SIXEL.
func IsCompatible() bool {
	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := os.Getenv("TERM_PROGRAM")

	// Terminals that advertise SIXEL support via $TERM
	sixelTerms := []string{"mlterm", "yaft", "foot", "contour", "domterm"}
	for _, t := range sixelTerms {
		if strings.Contains(term, t) {
			return true
		}
	}

	// Terminals that advertise support via $TERM_PROGRAM
	switch termProgram {
	case "iTerm.app", "WezTerm", "Tabby":
		return true
	}

	// xterm compiled with SIXEL support sets $XTERM_VERSION
	if strings.Contains(term, "xterm") && os.Getenv("XTERM_VERSION") != "" {
		return true
	}

	// mlterm sets $MLTERM when running
	if os.Getenv("MLTERM") != "" {
		return true
	}

	return false
}
