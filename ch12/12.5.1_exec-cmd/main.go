package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) == 1 {
		return
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	state := cmd.ProcessState
	fmt.Printf("%s\n", state.String())
	fmt.Printf("\tPid: %d\n", state.Pid())
	fmt.Printf("\tSystem Tyme: %v\n", state.SystemTime())
	fmt.Printf("\tUser Time: %v\n", state.UserTime())
}
