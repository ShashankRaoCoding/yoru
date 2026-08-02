package printpkg

import (
	"fmt"
	"yoru/sixel/print"
	"yoru/utils"
)

var subcommands = map[string]func([]string){
	"image": print.MainImage,
	"table": print.MainTable,
}

// Main dispatches to the appropriate print subcommand.
func Main(args []string) {
	if len(args) < 1 {
		utils.Error(fmt.Errorf("no subcommand provided; available subcommands: image, table"))
		return
	}

	subcmd := args[0]
	fn, ok := subcommands[subcmd]
	if !ok {
		utils.Error(fmt.Errorf("unknown subcommand %q; available subcommands: image, table", subcmd))
		return
	}

	fn(args[1:])
}
