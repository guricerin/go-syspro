package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Printf("Temp File Path: %s\n", filepath.Join(os.TempDir(), "temp.txt"))

	dir, name := filepath.Split(os.Getenv("HOME"))
	fmt.Printf("Dir: %s, Name: %s\n", dir, name)

	fragments := strings.Split(os.Getenv("HOME"), string(filepath.Separator))
	fmt.Printf("Fragments: %v\n", fragments)
}
