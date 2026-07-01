package make


import (
	"strings"
	"fmt" 
	"os/exec" 
	"os" 
	"yoru/globals"
)

func Init() {
	globals.Methods["mkfile"] = mkfile
}

func mkfile(filenames [] string) error {
	var err error 
	for _, name := range filenames {
		cmd := exec.Command("python", "-c", fmt.Sprintf(script, name)) 
		cmd.Stdout = os.Stdout 
		cmd.Stderr = os.Stderr 
		err = cmd.Start() 
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not create %s due to the error %s\n", name, err) 
			continue 
		}
		err = cmd.Wait() 
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not create %s due to the error %s\n", name, err) 
			continue 
		}
	}

	return nil 
}

var script string = `
file = open(%s, "w")
file.write("\n".join("" for i in range(45)))
file.close() 
`





































