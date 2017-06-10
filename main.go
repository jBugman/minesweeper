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

	for {
		success := bot.GameLoop()
		if !success {
			bot.StartGame()
		}
	}
}
