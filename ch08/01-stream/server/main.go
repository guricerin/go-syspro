package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Unixドメインソケット: カーネル内部で完結する高速なネットワークインターフェース
	// Unixドメインソケットはソケットファイルを作成し、これを介してプロセル間を通信させる
	path := filepath.Join(os.TempDir(), "unixdomainsocket-sample")
	// ソケットファイルが既存だとソケットが開けないので削除
	// なければ、それはそれで問題ない
	os.Remove(path)
	listener, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}
	// Closeしないとソケットファイルが残り続けてしまう
	defer listener.Close()

	fmt.Println("Server is running at " + path)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// リクエストを読み込む
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(dump))
			// レスポンスを書き込む
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
