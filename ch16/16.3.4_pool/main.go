package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var count int
	// 内部では要素はキャッシュでしかなく、他言語でいうところのWeakRefのコンテナ
	// -> GCが走ると保持しているデータが削除される。sync.Poolは消えては困るデータのコンテナには適さない。
	pool := sync.Pool{
		New: func() interface{} {
			count++
			return fmt.Sprintf("created: %d", count)
		},
	}

	// 追加した要素から受け取れる
	// プールが空だと新規作成
	pool.Put("manualy added: 1")
	pool.Put("manualy added: 2")
	fmt.Println(pool.Get())
	fmt.Println(pool.Get())
	fmt.Println(pool.Get()) // これは新規作成

	// GCを呼ぶと追加された要素が消える
	pool.Put("removed: 1")
	pool.Put("removed: 2")
	pool.Put("removed: 3")
	runtime.GC()
	fmt.Println(pool.Get())
}
