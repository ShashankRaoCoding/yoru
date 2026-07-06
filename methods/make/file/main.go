package file

import (
	// "strings"
	"fmt" 
	// "os/exec" 
	"os" 
	// "yoru/globals" 
)

var c string = ""

func init() {
	for range 45 {
		c += "\n" 
	}
}

func Main(filenames [] string) error {
	// var err error 
	for _, name := range filenames {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("There was an error creating the file", name+":", err.Error())
		}
		file.Write([]byte(c))
		file.Close() 
	}

	return nil 
}






































