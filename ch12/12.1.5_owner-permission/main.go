package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("ユーザID: %d\n", os.Getuid())
	fmt.Printf("グループID: %d\n", os.Getgid())
	fmt.Printf("実行ユーザ（実行ファイルの所有者）ID: %d\n", os.Geteuid())
	fmt.Printf("実行グループ（実行ファイルの所有グループ）ID: %d\n", os.Getegid())
}
