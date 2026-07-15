package mkpkg

import (
	"os"
	"yoru/utils"
)

func Main(args []string) {
	for _, path := range args {
		file, err := os.Create(path)
		defer file.Close()
		utils.Error(err)
		text := ""
		for range 45 {
			text += "\n"
		}
		file.Write([]byte(text))
	}
}
