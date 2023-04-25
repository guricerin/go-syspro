package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	listner, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running at localhost:8888")
	// 一度のアクセスで終わらないよう、無限ループさせる
	for {
		conn, err := listner.Accept()
		if err != nil {
			panic(err)
		}

		// 1リクエスト処理中に他のリクエストをAccept()できるように
		// goroutineを使って非同期にレスポンスを処理する
		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// HTTPリクエストのヘッダーやメソッドをparse
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(dump))

			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       io.NopCloser(strings.NewReader("Hello World\n")),
			}
			response.Write(conn)
			conn.Close()
		}()
	}
}
