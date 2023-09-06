package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

// 指定時間以内にSIGINT or SIGTERMを受信したら文字列を出力して終了する
// Ctrl + C
func main() {
	ctx := context.Background()
	sigctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	toctx, cancel2 := context.WithTimeout(ctx, time.Second*5)
	defer cancel2()

	fmt.Println("waiting...")
	select {
	case <-sigctx.Done():
		fmt.Println("signal recieved")
	case <-toctx.Done():
		fmt.Println("timeout")
	}
}
