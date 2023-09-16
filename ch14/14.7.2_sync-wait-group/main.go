package main

import (
	"fmt"
	"sync"
)

func main() {
	// チャネルと比較した場合、
	// ジョブ数が大量だったり可変個のときにWaitGroupのほうが使いやすい
	var wg sync.WaitGroup
	wg.Add(2) // goroutine起動前にジョブ数はあらかじめ決めておく

	go func() {
		fmt.Println("task 1")
		wg.Done()
	}()

	go func() {
		fmt.Println("task 2")
		wg.Done()
	}()

	wg.Wait() // すべてのジョブの終了を待機
	fmt.Println("finished")
}
