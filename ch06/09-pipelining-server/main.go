package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// 順番に従ってconnに書き出す（goroutineで実行する想定）
func writeToConn(sessionResponses chan chan *http.Response, conn net.Conn) {
	defer conn.Close()
	for sessionResponse := range sessionResponses {
		// 選択された仕事が終わるまで待機
		response := <-sessionResponse
		response.Write(conn)
		close(sessionResponse)
	}
}

// セッション内のリクエストを処理する
func handleRequest(request *http.Request, resultReciever chan *http.Response) {
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
	// content := "Hello World\n"
	content := string(dump) + "\n" // サーバ側のログは順番がバラバラで、クライアント側は順番通りか試す。
	// セッションを維持するためにKeep-Aliveでないといけない
	response := &http.Response{
		StatusCode:    200,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(len(content)),
		Body:          io.NopCloser(strings.NewReader(content)),
	}
	// 処理が終わったらチャネルに書き込み、
	// ブロックされていたwriteToConnの処理を再始動する
	resultReciever <- response
}

// 1セッションの処理を行う
func processSession(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	// セッション内のリクエストを順に処理するためのキューとしてのバッファ付きチャネル
	// その内部は送信データをためるバッファなしチャネル
	sessionResponses := make(chan chan *http.Response, 50)
	defer close(sessionResponses)
	// レスポンスを直列化してソケットに書き出す専用のgoroutine
	go writeToConn(sessionResponses, conn)
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		request, err := http.ReadRequest(reader)
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				fmt.Println("Timeout")
				break
			} else if err == io.EOF {
				break
			}
			panic(err)
		}
		// 受け取ったレスポンスをセッションのキューに入れる
		sessionResponse := make(chan *http.Response)
		sessionResponses <- sessionResponse
		// リクエストごとに非同期でレスポンスを実行
		go handleRequest(request, sessionResponse)
	}
}

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
		go processSession(conn)
	}
}
