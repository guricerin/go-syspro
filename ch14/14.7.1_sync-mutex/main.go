package main

import (
	"fmt"
	"sync"
	"time"
)

var id int

func generateId(mutex *sync.Mutex) int {
	// sync.Mutexを使うと、メモリを読み書きするコードに入るgoroutineがひとつに制限される
	// -> レースコンディションを回避可能
	mutex.Lock()
	// --> クリティカルセクション ここから
	id++
	result := id
	// <-- クリティカルセクション ここまで
	mutex.Unlock()
	return result
}

func main() {
	// mutex := new(sync.Mutex)とほぼ等価（←はmutexがポインタ型になる）
	var mutex sync.Mutex

	for i := 0; i < 100; i++ {
		go func() {
			fmt.Printf("id: %d\n", generateId(&mutex))
		}()
	}
	time.Sleep(time.Second)
}
