package main

import "os"

func Main(args []string) {
	for _, path := range args {
		file, err := os.Create(path)
		defer file.Close()
		globals.Handle(err)
		text := ""
		for range 45 {
			text += "\n"
		}
		file.Write([]byte(text))
	}
}
