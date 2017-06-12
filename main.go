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

	var success bool
	for !success {
		success = bot.GameLoop()
	}
}
