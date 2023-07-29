package main

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/google/shlex"
	"github.com/peterh/liner"
)

func parseCmd(cmdStr string) (cmd string, args []string, err error) {
	l := shlex.NewLexer(strings.NewReader(cmdStr))
	cmd, err = l.Next()
	if err != nil {
		return
	}
	for {
		token, err := l.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			break
		}

		args = append(args, token)
	}
	return
}

func main() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	for {
		if cmdline, err := line.Prompt("> "); err == nil {
			if cmdline == "" {
				continue
			}
			// コマンド処理
			cmd, args, err := parseCmd(cmdline)
			if err != nil {
				log.Print("Error lex line: ", err)
				continue
			}
			log.Print(cmd, args)
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
