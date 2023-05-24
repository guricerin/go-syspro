package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

func main() {
	conn, err := net.Dial("unix", filepath.Join(os.TempDir(), "unixdomainsocket-sample"))
	if err != nil {
		panic(err)
	}
	// ホスト名を指定しているが。接続先をファイル名で指定しているので、この情報はリクエストヘッダー以外では使われない
	request, err := http.NewRequest("get", "http://localhost:8888", nil)
	if err != nil {
		panic(err)
	}
	request.Write(conn)

	response, err := http.ReadResponse(bufio.NewReader(conn), request)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}
