package main

import (
	"log"
	"time"

	"./engine"
)

func main() {
	bot := engine.NewEngine()
	bot.SetClickDuration(15 * time.Millisecond)
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}

	maxRetries := 5
	var success bool
	for !success {
		bot.StartGame()

		success = bot.GameLoop()
		if success {
			break
		}

		maxRetries--
		if maxRetries > 0 {
			log.Println("Attempts remaining:", maxRetries)
		} else {
			break
		}
	}
}
