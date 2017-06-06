package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"./macos"
)

func saveImage(img image.Image) {
	outFile, _ := os.Create("test.png")
	png.Encode(outFile, img)
}

const targetAppTitle = "Minesweeper"

func main() {
	t := time.Now()
	winMeta, err := macos.FindWindow(targetAppTitle)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(winMeta)
	screenShot := macos.TakeScreenshot(winMeta.ID)
	log.Println(time.Since(t))

	saveImage(screenShot)
}
