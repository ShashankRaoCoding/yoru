package print

import (
	"fmt"

	"yoru/utils"
)

type runnerFunc func([]string) error

var subcommands = map[string]runnerFunc{
	"image": runImage,
	"table": runTable,
}

// Main dispatches to the appropriate sixel print subcommand.
func Main(args []string) {
	fn, remaining, err := resolveSubcommand(args)
	utils.Error(err)
	utils.Error(fn(remaining))
}

func resolveSubcommand(args []string) (runnerFunc, []string, error) {
	if len(args) < 1 {
		return nil, nil, fmt.Errorf("no subcommand provided; available subcommands: image, table")
	}

	subcmd := args[0]
	fn, ok := subcommands[subcmd]
	if !ok {
		return nil, nil, fmt.Errorf("unknown subcommand %q; available subcommands: image, table", subcmd)
	}

	return fn, args[1:], nil
}
