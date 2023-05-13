package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server is running at localhost:8888")
	// TCPと違い、Acceptは不要
	conn, err := net.ListenPacket("udp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// フレームサイズやバイナリ形式に従って読み込むので、バッファを用意
	buffer := make([]byte, 1500)
	for {
		// 送信内容 & 接続元の情報を受取る
		// ReadFrom()ではなくRead()だと、レスポンスを返す必要がある場合に対応不可能
		length, remoteAddress, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received from %v: %v\n", remoteAddress, string(buffer[:length]))
		_, err = conn.WriteTo([]byte("Hello from Server"), remoteAddress)
		if err != nil {
			panic(err)
		}
	}
}
