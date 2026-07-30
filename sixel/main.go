package sixelpkg

import (
	"fmt"
	"yoru/sixel/compatible"
	sixelprint "yoru/sixel/print"
	"yoru/utils"
)

var subcommands = map[string]func([]string){
	"compatible": compatible.Main,
	"print":      sixelprint.Main,
}

// Main dispatches to the appropriate sixel subcommand.
func Main(args []string) {
	if len(args) < 1 {
		utils.Error(fmt.Errorf("no subcommand provided; available subcommands: compatible, print"))
		return
	}

	subcmd := args[0]
	fn, ok := subcommands[subcmd]
	if !ok {
		utils.Error(fmt.Errorf("unknown subcommand %q; available subcommands: compatible, print", subcmd))
		return
	}

	fn(args[1:])
}
