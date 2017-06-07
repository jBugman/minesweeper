package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"./engine"
)

func saveImage(img image.Image) {
	outFile, _ := os.Create("test.png")
	png.Encode(outFile, img)
}

func main() {
	bot := engine.NewEngine()
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}

	bot.StartGame()
	img := bot.GrabScreen()
	saveImage(img)
	// Click in some random points
	bot.LeftClick(10, 10)
	bot.RightClick(0, 2)
	bot.RightClick(3, 2)
}
