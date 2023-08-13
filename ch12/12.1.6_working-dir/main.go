package main

import (
	"fmt"
	"os"
)

func main() {
	// 作業フォルダはプロセスの中で1つだけ設定可能
	// マルチスレッドでも、スレッド毎に別の作業フォルダを設定することは不可能
	wd, _ := os.Getwd()
	fmt.Println(wd)
}
