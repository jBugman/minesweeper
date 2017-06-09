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

	// bot.StartGame()
	bot.GameLoop()
	// Click in some random points
	// bot.LeftClick(6, 6)
	// bot.RightClick(0, 2)
	// bot.RightClick(3, 2)
	// bot.UpdateField()
	// bot.PrintField()
}
