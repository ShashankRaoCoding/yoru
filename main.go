package main

import (
	"yoru/methods" 
	"os" 
	"fmt" 
	// mk "yoru/make" 
	// "yoru/csv"
)

var Methods = make(map[string]func([]string)error) 

func init() {
	Methods["methods"] = methods.Main
	Methods["info"] = info
}

func main() {

	args := os.Args[1:] 

	method, args := args[0], args[1:] 
	
	m, ok := Methods[method]
	if ok == false {
		fmt.Printf("Function %s does not exist \n", method)
		os.Exit(1) 
	}

	err := m(os.Args[2:]) 

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1) 
	} 

	os.Exit(0) 
}

func info(args []string) error {
	fmt.Println("Options: ") 
	for name, _ := range Methods {
		fmt.Println(name) 
	}
	return nil 
}












































