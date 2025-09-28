package main

import (
	"os"
)

func main() {

	bot, err := NewBot(
		os.Getenv("MM_TOKEN"),
		os.Getenv("MM_SERVER"),
		os.Getenv("MM_TEAM"),
	)

	if err != nil {
		panic(err)
	}

	bot.Listen()
}
