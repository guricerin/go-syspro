package main

import (
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}
	err = req.Write(conn)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		panic(err)
	}
}
