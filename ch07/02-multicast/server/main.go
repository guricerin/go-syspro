package main

import (
	"fmt"
	"net"
	"time"
)

const interval = 3 * time.Second

// UDPのマルチキャストでは、クライアントがソケットを開いて待ち受け、サーバがデータを送信
func main() {
	fmt.Println("Start tick server at 224.0.0.1:9999")
	conn, err := net.Dial("udp", "224.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	start := time.Now()
	wait := start.Truncate(interval).Add(interval).Sub(start)
	time.Sleep(wait)
	ticker := time.Tick(interval)
	// interval (3秒) 間隔でfor文が回る
	for now := range ticker {
		conn.Write([]byte(now.String()))
		fmt.Println("Tick: ", now.String())
	}
}
