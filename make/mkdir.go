package make


import (
	// "strings"
	"fmt" 
	"os/exec" 
	"os" 
	"yoru/globals" 
)

func Init() {
	globals.Methods["mkdir"] = mkdir
}

func mkdir(dirnames [] string) error {
	var err error
	cmd := exec.Command("mkdir", "-p", dirnames...) 
	cmd.Stdout = os.Stdout 
	cmd.Stderr = os.Stderr 
	err = cmd.Start() 
	if err != nil {
		return err
	}
	err = cmd.Wait() 
	if err != nil {
		return err 
	}

	return nil 
}







































