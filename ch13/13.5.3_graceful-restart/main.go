package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lestrrat-go/server-starter/listener"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	// Server::starterからもらったソケットを確認
	listners, err := listener.ListenAll()
	if err != nil {
		panic(err)
	}
	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "server pid: %d %v\n", os.Getpid(), os.Environ())
		}),
	}
	// サーバをgoroutineで起動
	go server.Serve(listners[0])

	// SIGTERMを受信したら終了させる
	<-signals
	server.Shutdown(context.Background())
}
