package engine

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"../macos"
	"../macos/keycode"
)

const (
	tileSize     = 32 // 64px on retina
	headerHeight = 22
	footerHeight = 31
)

type engine struct {
	x, y          int
	width, height uint
	windowID      int
	field         [][]uint8
}

// Engine provides public interface
type Engine interface {
	Start() error
	GrabScreen() image.Image
	StartGame()
	LeftClick(x, y int)
	RightClick(x, y int)
	PrintField()
	UpdateField()
}

// NewEngine creates engine instance
func NewEngine() Engine {
	return &engine{}
}

func (e *engine) Start() error {
	const appTitle = "Minesweeper"
	winMeta, err := macos.FindWindow(appTitle)
	log.Printf("%+v\n", winMeta)
	if err != nil {
		return err
	}

	e.windowID = winMeta.ID
	e.x = winMeta.Bounds.X()
	e.y = winMeta.Bounds.Y() + headerHeight
	e.width = winMeta.Bounds.Width() / tileSize
	e.height = (winMeta.Bounds.Height() - headerHeight - footerHeight) / tileSize

	// single-allocation method
	e.field = make([][]uint8, e.height)
	cells := make([]uint8, e.width*e.height)
	for i := range e.field {
		e.field[i], cells = cells[:e.width], cells[e.width:]
	}

	log.Printf("%dx%d", e.width, e.height)

	macos.ActivateWindow(winMeta.OwnerPID)
	return nil
}

func (e *engine) GrabScreen() image.Image {
	img := macos.TakeScreenshot(e.windowID)
	cropped := img.SubImage(rect(0, headerHeight, e.width*tileSize, headerHeight+e.height*tileSize))
	return cropped
}

func (e engine) StartGame() {
	macos.KeyPressWithModifier(keycode.KeyN, keycode.KeyCommand)
}

func (e engine) tileCenterX(x int) int {
	return e.x + x*tileSize + tileSize/2
}

func (e engine) tileCenterY(y int) int {
	return e.y + y*tileSize + tileSize/2
}

func (e engine) LeftClick(x, y int) {
	macos.LeftClick(e.tileCenterX(x), e.tileCenterY(y))
}

func (e engine) RightClick(x, y int) {
	macos.RightClick(e.tileCenterX(x), e.tileCenterY(y))
}

func (e engine) PrintField() {
	for _, line := range e.field {
		log.Println(line)
	}
}

func rect(x0, y0, x1, y1 uint) image.Rectangle {
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func (e engine) UpdateField() {
	img := e.GrabScreen().(*image.RGBA)
	saveImage("debug/field.png", img)
	var i, j uint
	for i = 0; i < e.width; i++ {
		for j = 0; j < e.height; j++ {
			tile := img.SubImage(rect(i*tileSize, j*tileSize+headerHeight, (i+1)*tileSize, (j+1)*tileSize+headerHeight))
			saveImage(fmt.Sprintf("debug/test_%d_%d.png", i, j), tile)
		}
	}
}

func saveImage(filename string, img image.Image) {
	outFile, _ := os.Create(filename)
	defer outFile.Close()
	png.Encode(outFile, img)
}
