package filegeneration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fanialfi/fan-out-fan-in/lib"
)

func GenerateFileSequential() {
	err := os.RemoveAll(lib.TempPath)
	if err != nil {
		log.Printf("ERROR remove all directory and file : %s\n", err.Error())
	}

	err = os.MkdirAll(lib.TempPath, os.ModePerm)
	if err != nil {
		log.Printf("ERROR create directory : %s\n", err.Error())
	}

	for i := 0; i < lib.TotalFile; i++ {
		filename := filepath.Join(lib.TempPath, fmt.Sprintf("file-%d.txt", i))
		content := lib.RandomString(lib.ContentLength)
		err := os.WriteFile(filename, []byte(content), os.ModePerm)
		if err != nil {
			log.Printf("ERROR write data to file : %s\n", err.Error())
		}

		log.Println(i, "file created")
	}

	log.Printf("%d file was created\n", lib.TotalFile)
}
