package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"./macos"
	"./macos/keycode"
)

func saveImage(img image.Image) {
	outFile, _ := os.Create("test.png")
	png.Encode(outFile, img)
}

const targetAppTitle = "Minesweeper"

func main() {
	winMeta, err := macos.FindWindow(targetAppTitle)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(winMeta)

	macos.ActivateWindow(winMeta.OwnerPID)

	t := time.Now()
	screenShot := macos.TakeScreenshot(winMeta.ID)
	log.Println(time.Since(t))

	// Start new game
	macos.KeyPressWithModifier(keycode.KeyN, keycode.KeyCommand)
	// Click in some random points
	macos.LeftClick(winMeta.Bounds.X()+100, winMeta.Bounds.Y()+200)
	macos.RightClick(winMeta.Bounds.X()+250, winMeta.Bounds.Y()+250)
	macos.LeftClick(winMeta.Bounds.X()+300, winMeta.Bounds.Y()+500)

	_ = screenShot
	// saveImage(screenShot)
}
