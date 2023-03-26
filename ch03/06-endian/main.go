package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	// 32bitのビッグエンディアンのデータ 10000(0x2710)
	data := []byte{0x00, 0x00, 0x27, 0x10}
	var i int32
	// エンディアン変換
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Printf("data: %d\n", i)
}
