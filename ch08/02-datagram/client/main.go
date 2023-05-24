package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	// 受信用のソケットファイル
	// サーバのReadFrom()で返信可能なアドレスが得られる
	clientPath := filepath.Join(os.TempDir(), "unixdomainsocket-client")
	os.Remove(clientPath)
	conn, err := net.ListenPacket("unixgram", clientPath)
	if err != nil {
		panic(err)
	}
	// 送信先のアドレス
	unixServerAddr, err := net.ResolveUnixAddr("unixgram", filepath.Join(os.TempDir(), "unixdomainsocket-server"))
	if err != nil {
		panic(err)
	}
	var serverAddr net.Addr = unixServerAddr
	defer conn.Close()

	log.Println("Sending to server")
	_, err = conn.WriteTo([]byte("Hello from Client"), serverAddr)
	if err != nil {
		panic(err)
	}

	log.Println("Receiving from server")
	buffer := make([]byte, 1500)
	length, _, err := conn.ReadFrom(buffer)
	if err != nil {
		panic(err)
	}
	log.Printf("Received: %s\n", string(buffer[:length]))
}
