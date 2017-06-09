package engine

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/jBugman/imghash"

	"../macos"
	"../macos/keycode"
)

const (
	tileSize     = 32 // 64px on retina
	headerHeight = 22
	footerHeight = 31
)

// Tile represents possible tile values
type Tile uint8

const (
	OpenSpace Tile = 0
	Flag      Tile = 32
	Unknown   Tile = 255
)

var tileHashes = map[uint64]Tile{
	0x0000000000000000: OpenSpace, // OpenSpace OR Unknown tile
	0xFFC3C3E7E7E3E3FF: 1,
	0xFFC3E3E7CFC3E3FF: 2,
	0xFFE3C3C7E7C3E3FF: 3,
	0xFFEFC3C3E3E7EFFF: 4,
	// TODO other numbers
	0xFFF7F7C3C381F3FF: Flag,
}

type engine struct {
	x, y          int
	width, height uint
	windowID      int
	field         [][]Tile
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
	e.field = make([][]Tile, e.height)
	cells := make([]Tile, e.width*e.height)
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

func (e *engine) UpdateField() {
	img := e.GrabScreen().(*image.RGBA)
	saveImage("debug/field.png", img)
	var i, j uint
	for i = 0; i < e.width; i++ {
		for j = 0; j < e.height; j++ {
			tile := img.SubImage(rect(i*tileSize, j*tileSize+headerHeight, (i+1)*tileSize, (j+1)*tileSize+headerHeight))
			tileValue := recognizeTile(tile)
			e.field[j][i] = tileValue
		}
	}
}

func recognizeTile(tile image.Image) Tile {
	hash := imghash.Average(tile)
	value, ok := tileHashes[hash]
	if !ok {
		saveImage(fmt.Sprintf("debug/error_%X.png", hash), tile)
		log.Fatalf("Unknown hash: %X\n", hash)
	}
	if value == 0 {
		// tile is a subimage, so we need its offset
		coords := tile.Bounds().Min
		col := tile.(*image.RGBA).RGBAAt(coords.X+1, coords.Y+1)

		avg := (uint(col.R) + uint(col.G) + uint(col.B)) / 3
		if avg < 180 { // if it is more purple than white
			value = Unknown
		}
	}
	return value
}

func saveImage(filename string, img image.Image) {
	outFile, _ := os.Create(filename)
	defer outFile.Close()
	png.Encode(outFile, img)
}
