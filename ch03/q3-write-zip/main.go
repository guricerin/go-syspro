package main

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

func main() {
	zipFile, err := os.Create("out.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	file, err := zipWriter.Create("newFile.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(file, strings.NewReader("aaa\nbbbbb"))
}
