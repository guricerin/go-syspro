package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// パスをそのままクリーンにする
	fmt.Println(filepath.Clean("./path/filepath/../path.go"))
	// path/path.go

	// 絶対パスに整形
	abspath, _ := filepath.Abs("path/filepath/path_unix.go")
	fmt.Println(abspath)

	// 相対パスに整形
	relpath, _ := filepath.Rel("/usr/local/go/src", "/usr/local/go/src/path/filepath/path.go")
	fmt.Println(relpath)

	// 環境変数の展開
	path := os.ExpandEnv("${HOME}/src/github.com/shibukawa/tobubus")
	fmt.Println(path)

	// $HOME
	fmt.Println(os.UserHomeDir())
}
