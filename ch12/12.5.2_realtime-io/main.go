package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	count := exec.Command("./count.out")
	stdout, _ := count.StdoutPipe()
	go func() {
		// パイプで繋いできた子プロセスの標準出力、このプロセス（親）にリダイレクト
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("(stdout) %s\n", scanner.Text())
		}
	}()
	err := count.Run()
	if err != nil {
		panic(err)
	}
}
