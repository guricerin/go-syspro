package main

import (
	"bufio"
	"fmt"
	"strings"
)

var source = `1行目
2行目
3行目`

// func main() {
// 	// 分割文字は残る
// 	reader := bufio.NewReader(strings.NewReader(source))
// 	for {
// 		line, err := reader.ReadString('\n')
// 		fmt.Printf("%#v\n", line)
// 		if err == io.EOF {
// 			break
// 		}
// 	}
// }

func main() {
	// 分割文字は残らない
	scanner := bufio.NewScanner(strings.NewReader(source))
	// デフォルトは改行区切りなのを単語区切りにする
	// scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Printf("%#v\n", scanner.Text())
	}
}
