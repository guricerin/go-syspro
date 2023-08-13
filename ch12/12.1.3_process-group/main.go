package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	// プロセスグループ: パイプで繋がったプロセス群
	// セッショングループ: 同じターミナルから起動したプロセス群

	// linux_amd64 に syscall.Getsid はなかった...
	// sid, _ := syscall.Getsid(os.Getpid())
	sid := 0
	fmt.Fprintf(os.Stderr, "プロセスグループID: %d, セッショングループID: %d\n", syscall.Getpgrp(), sid)
}
