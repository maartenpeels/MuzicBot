package main

import (
	"muzicBot/bot"
	"os"
)

func main() {
	err := bot.LoadEnv()
	if err != nil {
		os.Exit(1)
	}

	bot.Init()
}
