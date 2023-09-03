package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("top")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Process.Signal(os.Interrupt)
	fmt.Println("Finished")
}
