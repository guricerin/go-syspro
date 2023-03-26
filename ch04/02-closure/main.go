package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start closure")
	go func() {
		fmt.Println("closure is running")
		time.Sleep(time.Second)
		fmt.Println("closure is finished")
	}()
	time.Sleep(2 * time.Second)
}
