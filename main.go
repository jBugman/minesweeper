package main

import (
	"log"

	"./engine"
)

func main() {
	bot := engine.NewEngine()
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
