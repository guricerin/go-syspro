package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Page Size: %d Byte\n", os.Getpagesize())
}
