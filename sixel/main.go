package sixelpkg

import (
	"fmt"
	"yoru/sixel/compatible"
	"yoru/utils"
)

var subcommands = map[string]func([]string){
	"compatible": compatible.Main,
}

// Main dispatches to the appropriate sixel subcommand.
func Main(args []string) {
	if len(args) < 1 {
		utils.Error(fmt.Errorf("no subcommand provided; available subcommands: compatible"))
		return
	}

	subcmd := args[0]
	fn, ok := subcommands[subcmd]
	if !ok {
		utils.Error(fmt.Errorf("unknown subcommand %q; available subcommands: compatible", subcmd))
		return
	}

	fn(args[1:])
}
