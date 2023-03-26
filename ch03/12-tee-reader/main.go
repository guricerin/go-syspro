package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	var buffer bytes.Buffer
	reader := bytes.NewBufferString("Example of io.TeeReader\n")
	teeReader := io.TeeReader(reader, &buffer)
	_, _ = io.ReadAll(teeReader)

	fmt.Println(buffer.String())
}
