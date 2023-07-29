package main

import (
	"errors"
	"io"
	"log"

	"github.com/peterh/liner"
)

func main() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	for {
		if cmd, err := line.Prompt("> "); err == nil {
			if cmd == "" {
				continue
			}
			// コマンド処理
			log.Print(cmd)
		} else if errors.Is(err, io.EOF) {
			break
		} else if err == liner.ErrPromptAborted {
			log.Print("Aborted")
			break
		} else {
			log.Print("Error reading line: ", err)
		}
	}
}
