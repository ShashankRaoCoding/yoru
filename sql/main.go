package sqlpkg

import (
	"io"
	"os"
	"yoru/globals"
	"yoru/utils"
)

func Main(args []string) {
	db, err := ImportCSVToSQL(os.Stdin, globals.TEMP, "table")
	utils.Error(err)
	
	query, err := io.ReadAll(os.NewFile(3, "QUERY"))
	utils.Error(err)
	
	err = QueryToCSV(db, string(query), os.Stdout)
	utils.Error(err)
}
