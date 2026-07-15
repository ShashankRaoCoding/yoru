package globals

import (
	"fmt"
	"os"
	"yoru/utils"
)

var TEMP string

func init() {
	temp, err := utils.GetBinaryPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error getting the binary path, there may be unexpected behaviour. TEMP files shall be stored in the working dir\n")
	} else {
		TEMP = temp
		return
	}
	TEMP = "./temp"
}
