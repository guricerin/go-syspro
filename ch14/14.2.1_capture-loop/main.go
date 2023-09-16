package main

import (
	"fmt"
	"time"
)

func main() {
	tasks := []string{
		"cmake ..",
		"cmake . --build Release",
		"cpack",
	}
	for _, task := range tasks {
		go func() {
			// goroutineがｈ起動するときにはループが回り切って
			// 全てのtaskがtasks[last_index]になってしまう
			// 対策: 関数の引数経由で明示的にデータのコピーが行われるようにする
			fmt.Println(task)
		}()
	}
	time.Sleep(time.Second)
}
