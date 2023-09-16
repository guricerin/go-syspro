package main

import (
	"fmt"
	"time"
)

func sub1(c int) {
	fmt.Println("share by arguments: ", c*c)
}
func main() {
	// 引数渡し
	go sub1(10)

	// クロージャのキャプチャ渡し
	c := 20
	go func() {
		fmt.Println("share by capture: ", c*c)
	}()
	// 本来はチャネルや sync.WaitGroup などの「作業が完了した」ことをきちんと取り扱える仕組みを使って待ち合わせ処理を書くほうが望ましい
	time.Sleep(time.Second)
}
