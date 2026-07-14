package main

import (
	"os"
	"io" 
	"yoru/globals"
)

func main() {
	db := ImportCSVToSQL(os.Stdin, globals.TEMP, "table")
	query, err := io.ReadAll(os.NewFile(3, "QUERY"))
	utils.Error(err) 
	err := QueryToCSV(db, string(query), os.Stdout) 
	utils.Error(err) 

}

