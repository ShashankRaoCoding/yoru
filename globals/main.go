
package globals

import (
	"yoru/utils"
	"os"
)

func init() {
	TEMP, err  := utils.GetBinaryPath() 
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error getting the binary path, there may be unexpected behaviour. TEMP files shall be stored in the working dir\n")
	}
	TEMP = "./temp" 
}

var TEMP string 

