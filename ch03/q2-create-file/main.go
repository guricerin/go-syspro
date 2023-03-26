package main

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
)

func main() {
	file, err := os.Create("sample.bin")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	_, err = rand.Reader.Read(buffer)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewBuffer(buffer)
	_, err = io.Copy(file, reader)
	if err != nil {
		panic(err)
	}
}
