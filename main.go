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

	macos.ActivateWindow(winMeta.OwnerPID)
	screenShot := macos.TakeScreenshot(winMeta.ID)

	log.Println(time.Since(t))

	macos.RightClick(winMeta.Bounds.X()+250, winMeta.Bounds.Y()+250)
	macos.LeftClick(winMeta.Bounds.X()+100, winMeta.Bounds.Y()+200)
	macos.LeftClick(winMeta.Bounds.X()+300, winMeta.Bounds.Y()+500)

	saveImage(screenShot)
}
