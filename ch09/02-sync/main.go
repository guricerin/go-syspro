package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	f, _ := os.Create("file.txt")
	a := time.Now()
	f.Write([]byte("緑の怪獣")) // この段階ではバッファメモリに書き込んだだけ
	b := time.Now()
	f.Sync() // 確実にストレージに書き込む
	c := time.Now()
	f.Close()
	d := time.Now()
	fmt.Printf("Write: %v\n", b.Sub(a))
	fmt.Printf("Sync: %v\n", c.Sub(b))
	fmt.Printf("Close: %v\n", d.Sub(c))
}
