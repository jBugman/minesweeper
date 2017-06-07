package engine

import (
	"image"
	"log"

	"../macos"
	"../macos/keycode"
)

const (
	tileSize     = 32 // 64px on retina
	headerHeight = 22
)

type engine struct {
	x, y          int
	width, height uint
	windowID      int
}

// Engine provides public interface
type Engine interface {
	Start() error
	GrabScreen() image.Image
	StartGame()
	LeftClick(x, y int)
	RightClick(x, y int)
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
	e.width = winMeta.Bounds.Width()
	e.height = winMeta.Bounds.Height()

	macos.ActivateWindow(winMeta.OwnerPID)
	return nil
}

func (e *engine) GrabScreen() image.Image {
	img := macos.TakeScreenshot(e.windowID)
	return img
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
