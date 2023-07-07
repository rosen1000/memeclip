package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.design/x/clipboard"
	"golang.org/x/exp/slices"

	"github.com/MarinX/keylogger"
	"github.com/go-vgo/robotgo"
)

var memes []string

func main() {
	go run()

	CatchErr(clipboard.Init())
	ch := clipboard.Watch(context.Background(), clipboard.FmtText)
	for data := range ch {
		fmt.Println(string(data))
		memes = append(memes, string(data))
	}
}

func run() {
	k, err := keylogger.New(keylogger.FindKeyboardDevice())
	CatchErr(err)

	var history []uint16

	defer k.Close()

	events := k.Read()
	for event := range events {
		if event.KeyPress() {
			history = append(history, event.Code)
			check := []uint16{64}
			if len(history) >= len(check) {
				if slices.Compare(history[len(history)-len(check):], check) == 0 {
					for _, meme := range memes {
						robotgo.TypeStr(meme)
						robotgo.KeyTap("enter")
						time.Sleep(time.Second)
					}
					memes = []string{}
					history = []uint16{}
					fmt.Println("Cleared clipboard")
				}
			}
			if event.Code == 1 {
				os.Exit(0)
			}
		}
	}
}

func CatchErr(err error) {
	if err != nil {
		panic(err)
	}
}
