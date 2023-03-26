package main

import (
	"flag"
	"io"
	"os"
)

func main() {
	oldFilePath := flag.String("i", "", "old file path")
	newFilePath := flag.String("o", "", "new file path")
	flag.Parse()

	oldFile, err := os.Open(*oldFilePath)
	if err != nil {
		panic(err)
	}
	defer oldFile.Close()
	newFile, err := os.Create(*newFilePath)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	io.Copy(newFile, oldFile)
}
