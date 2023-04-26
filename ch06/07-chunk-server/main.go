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

// 青空文庫: ごんぎつねより
// https://www.aozora.gr.jp/cards/000121/card628.html
var contents = []string{
	" これは、私わたしが小さいときに、村の茂平もへいというおじいさんからきいたお話です。",
	" むかしは、私たちの村のちかくの、中山なかやまというところに小さなお城があって、",
	" 中山さまというおとのさまが、おられたそうです。",
	" その中山から、少しはなれた山の中に、「ごん狐ぎつね」という狐がいました。",
	" ごんは、一人ひとりぼっちの小狐で、しだの一ぱいしげった森の中に穴をほって住んでいました。",
	" そして、夜でも昼でも、あたりの村へ出てきて、いたずらばかりしました。",
}

// 1セッションの処理を行う
func processSession(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	defer conn.Close()

	// Accept後のソケットで何度も応答を返すためにループ
	for {
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			if err == io.EOF {
				break
			}
			// 上記以外はエラーにする
			panic(err)
		}

		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		// チャンク形式ではヘッダーに送信データのサイズを書かない
		// だが、http.Responseはファイルサイズ指定がないと Connection: close を送信してしまうので
		// fmt.Fprintfで直接HTTPレスポンスを書き込む
		fmt.Fprintf(conn, strings.Join([]string{
			"HTTP/1.1 200 OK",
			"Content-Type: text/plain",
			"Transfer-Encoding: chunked",
			"",
			"",
		}, "\r\n"))
		for _, content := range contents {
			bytes := []byte(content)
			fmt.Fprintf(conn, "%x\r\n%s\r\n", len(bytes), content)
		}
		// チャンクの送信完了
		fmt.Fprintf(conn, "0\r\n\r\n")
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
