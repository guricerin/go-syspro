package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	sec := flag.Int("s", 1, "sec")
	flag.Parse()

	ch := time.After(time.Duration(*sec) * time.Second)
	done := make(chan bool)

	fmt.Println(time.Now())
	fmt.Println("start")
	go func() {
		for {
			select {
			case data := <-ch:
				fmt.Println(data)
				done <- true
				return
			default:
				break
			}
		}
	}()
	<-done
	fmt.Println("finished")
}
