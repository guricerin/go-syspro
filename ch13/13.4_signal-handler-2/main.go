package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SIGHUPなど繰り返し受け取るシグナルを扱うのは、こちらの書き方が便利
func main() {
	// サイズが1より大きいチャネルを作成
	signals := make(chan os.Signal, 1)
	// SIGINTを受け取る
	signal.Notify(signals, syscall.SIGINT)

	// 待機
	fmt.Println("Waiting SIGINT (Ctrl + C)")
	<-signals
	fmt.Println("SIGINT arrived")
}
