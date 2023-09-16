package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mutex sync.Mutex
	// 先行ジョブが終了してから待機中のすべてのジョブに通知を行う
	// チャネル経由でやりとりするデータが先行ジョブの終了フラグ以外になく、後発ジョブが大量な場合に、チャネルより便利
	// また、チャネルの場合は一度しか通知できない
	cond := sync.NewCond(&mutex)

	for _, name := range []string{"A", "B", "C"} {
		go func(n string) {
			cond.L.Lock()
			defer cond.L.Unlock()
			cond.Wait() // BroadCast()が実行されるまで待機
			fmt.Println(n)
		}(name)
	}

	fmt.Println("ready...")
	time.Sleep(time.Second) // 先行ジョブのつもり
	fmt.Println("go")
	cond.Broadcast()
	time.Sleep(time.Second)
}
