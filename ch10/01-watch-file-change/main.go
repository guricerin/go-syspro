package main

import (
	"log"

	"gopkg.in/fsnotify/fsnotify.v1"
)

func main() {
	counter := 0
	// パッシブ方式
	// 監視したいファイルをOS側に通知、ファイルが変更されたら教えてもらう
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events: // ブロッキングしてイベントを待機
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
					counter++
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					counter++
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file", event.Name)
					counter++
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("renamed file", event.Name)
					counter++
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("chmod file", event.Name)
					counter++
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
			if counter > 3 {
				done <- true
			}
		}
	}()

	// 監視対象フォルダを追加
	err = watcher.Add(".")
	if err != nil {
		// panic()でないのはなぜ？fatalだとdeferされない
		// ref: https://devlights.hatenablog.com/entry/2022/06/07/073000
		log.Fatal(err)
	}
	<-done
}
