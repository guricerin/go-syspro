package main

import (
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	second := 5
	duration := time.Duration(second)
	fmt.Printf("Accept Ctrl + C for %d seconds\n", second)
	time.Sleep(time.Second * duration)

	// 可変長引数で任意数のシグナルを無視できる
	// シグナルを受け取る訳では無いのでチャネルは不要
	signal.Ignore(syscall.SIGINT, syscall.SIGHUP)

	fmt.Printf("Ignore Ctrl + C for %d seconds\n", second)
	time.Sleep(time.Second * duration)
}
