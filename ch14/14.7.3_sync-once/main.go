package main

import (
	"fmt"
	"sync"
)

func initialize() {
	fmt.Println("init()より後に遅延させたい初期化処理")
}

var once sync.Once

func main() {
	// 最初の一回しか実行されない
	once.Do(initialize)
	once.Do(initialize)
	once.Do(initialize)
}
