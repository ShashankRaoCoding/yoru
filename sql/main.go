package sqlpkg

import (
	"fmt"
	"os"
	"yoru/globals"
	"yoru/utils"
)

func Main(args []string) {
	if len(args) < 1 {
		utils.Error(fmt.Errorf("usage: sql <query>"))
		return
	}

	query := args[0]

	db, err := ImportCSVToSQL(os.Stdin, globals.TEMP, "table")
	utils.Error(err)
	
	err = QueryToCSV(db, query, os.Stdout)
	utils.Error(err)
}
