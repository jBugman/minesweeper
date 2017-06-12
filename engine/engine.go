package engine

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"

	"github.com/jBugman/imghash"

	"../macos"
	"../macos/keycode"
)

const (
	tileSize        = 32 // 64px on retina
	headerHeight    = 22
	footerHeight    = 31
	restartMessageW = 214
	restartMessageH = 32
)
const (
	messageMask   ImageHash = 0x0300007E7E000000
	zeroBombsHash ImageHash = 0xFF8F2737373787CF
)

// Tile represents possible tile values
type Tile uint8

// ImageHash represents average image hash
type ImageHash uint64

// Special Tile values
const (
	OpenSpace Tile = 0
	Flag      Tile = 32
	Bomb      Tile = 64
	Unknown   Tile = 255
)

func (t Tile) String() string {
	switch t {
	case 0:
		return "🆓"
	case 1:
		return "1️⃣"
	case 2:
		return "2️⃣"
	case 3:
		return "3️⃣"
	case 4:
		return "4️⃣"
	case 5:
		return "5️⃣"
	case 6:
		return "6️⃣"
	case 7:
		return "7️⃣"
	case 8:
		return "8️⃣"
	case Flag:
		return "🚩"
	case Bomb:
		return "💣"
	case Unknown:
		return "❔"
	default:
		return string(t)
	}
}

var tileHashes = map[ImageHash]Tile{
	0x0000000000000000: OpenSpace, // OpenSpace OR Unknown tile
	0xFFC3C3E7E7E3E3FF: 1,
	0xFFC3E3E7CFC3E3FF: 2,
	0xFFE3C3C7E7C3E3FF: 3,
	0xFFEFC3C3E3E7EFFF: 4,
	0xFFE3C3CFE3E3E3FF: 5,
	0xFFE7C3C3E3E3E7FF: 6,
	0xFFF3F7E7E7C3C3FF: 7,
	0xFFE7C3C3E3C3E7FF: 8,
	0xFFF7F7C3C381F3FF: Flag,
	0xFFFFC3C3C3C3FFFF: Bomb,
}

type engine struct {
	x, y          int
	width, height uint
	windowID      int
	field         [][]Tile
	timerHash     ImageHash
	bombCountHash ImageHash
}

// Engine provides public interface
type Engine interface {
	Start() error
	GrabScreen() image.Image
	StartGame()
	LeftClick(x, y int)
	RightClick(x, y int)
	PrintField()
	UpdateField() error
	GameLoop() bool
	ClickRandomUnknown() bool
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
	cropped := img.SubImage(rect(0, headerHeight, e.width*tileSize, headerHeight+e.height*tileSize+footerHeight))
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
		log.Println(tilesString(line))
	}
}

func tilesString(tiles []Tile) string {
	var buf bytes.Buffer
	for _, tile := range tiles {
		buf.WriteString(fmt.Sprintf("%s ", tile))
	}
	return buf.String()
}

func rect(x0, y0, x1, y1 uint) image.Rectangle {
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func (e *engine) UpdateField() error {
	img := e.GrabScreen().(*image.RGBA)
	saveImage("debug/field.png", img)

	const (
		topMargin   = 9
		rightMargin = 20
		squareSize  = 16
	)
	bombs := img.SubImage(rect(
		e.width*tileSize-rightMargin-squareSize,
		e.height*tileSize+headerHeight+topMargin,
		e.width*tileSize-rightMargin,
		e.height*tileSize+headerHeight+topMargin+squareSize,
	))
	e.bombCountHash = ImageHash(imghash.Average(bombs))
	// saveImage("debug/bombs.png", bombs)
	// log.Printf("bomb hash: %X\n", e.bombCountHash)

	var i, j uint
	var err error
	var tileValue Tile
	for i = 0; i < e.width; i++ {
		for j = 0; j < e.height; j++ {
			tile := img.SubImage(rect(i*tileSize, j*tileSize+headerHeight, (i+1)*tileSize, (j+1)*tileSize+headerHeight))
			tileValue, err = e.recognizeTile(tile)
			if err != nil {
				return err
			}
			e.field[j][i] = tileValue
		}
	}
	return nil
}

func (e engine) recognizeTile(tile image.Image) (Tile, error) {
	hash := ImageHash(imghash.Average(tile))
	value, ok := tileHashes[hash]
	if !ok {
		saveImage(fmt.Sprintf("debug/error_%X.png", hash), tile)
		return Unknown, fmt.Errorf("Unknown hash: %X\n", hash)
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
	return value, nil
}

func saveImage(filename string, img image.Image) {
	outFile, _ := os.Create(filename)
	defer outFile.Close()
	png.Encode(outFile, img)
}

// GameLoop handles game logic and communication
func (e *engine) GameLoop() bool {
	var x, y int
	var didSomething bool
	var err error
	for {
		didSomething = false
		err = e.UpdateField()
		if e.bombCountHash == zeroBombsHash {
			log.Println("🤔 Should be victory but some tiles may remain")
			// Clicking on ramaining unknowns
			for y = 0; y < int(e.height); y++ {
				for x = 0; x < int(e.width); x++ {
					t := e.field[y][x]
					if t == Unknown {
						log.Println("Clicking on", x, y)
						e.LeftClick(x, y)
					}
				}
			}
			log.Println("🎉 Victory!")
			return true
		} else if err != nil {
			// log.Fatal(err) // Assume that we know all useful hashes
			log.Println("💣 Boom!")
			return false
		}
		e.PrintField()
		for y = 0; y < int(e.height) && !didSomething; y++ {
			for x = 0; x < int(e.width) && !didSomething; x++ {
				didSomething, err = e.processTile(x, y)
				if err != nil {
					log.Println(err)
					return false
				}
			}
		}
		if !didSomething {
			log.Println("🌀 Cannot decide what to do..")
			if !e.ClickRandomUnknown() {
				return false
			}
		}
	}
}

func (e *engine) processTile(x, y int) (bool, error) {
	tile := e.field[y][x]
	if tile == Bomb {
		return false, errors.New("😱 Bombs on the field! Starting again")
	}
	// log.Println(x, y, tile)
	if tile < 1 || tile > 8 {
		return false, nil
	}
	_, coords, unknownCount, flagCount := e.getNeighbours(x, y)
	// log.Println(tilesString(tiles))
	// Marking flags
	if unknownCount == int(tile)-flagCount {
		for k := 0; k < len(coords); k++ {
			c := coords[k]
			t := e.field[c.Y][c.X]
			if t == Unknown {
				e.field[c.Y][c.X] = Flag
				flagCount++
				log.Println("Setting flag at", c.X, c.Y)
				e.RightClick(c.X, c.Y)
				return true, nil
			}
		}
	}
	// Clicking on safe unknowns
	if unknownCount > 0 && int(tile) == flagCount {
		for k := 0; k < len(coords); k++ {
			c := coords[k]
			t := e.field[c.Y][c.X]
			if t == Unknown {
				log.Println("Clicking on", c.X, c.Y)
				e.LeftClick(c.X, c.Y)
				return true, nil
			}
		}
	}
	return false, nil
}

func (e *engine) ClickRandomUnknown() bool {
	var unknownCount int
	for y := 0; y < int(e.height); y++ {
		for x := 0; x < int(e.width); x++ {
			tile := e.field[y][x]
			if tile == Unknown {
				unknownCount++
			}
		}
	}
	if unknownCount == 0 {
		return false
	}
	randomIndex := rand.Intn(unknownCount)
	unknownCount = 0
	for y := 0; y < int(e.height); y++ {
		for x := 0; x < int(e.width); x++ {
			tile := e.field[y][x]
			if tile == Unknown {
				if unknownCount == randomIndex {
					log.Println("❗️ Randomly clicking on", x, y)
					e.LeftClick(x, y)
					return true
				}
				unknownCount++
			}
		}
	}
	return false
}

func (e engine) getNeighbours(x0, y0 int) ([]Tile, []image.Point, int, int) {
	var tiles = []Tile{}
	var coords = []image.Point{}
	var unknownCount, flagCount int
	var tile Tile
	for y := max(0, y0-1); y < min(int(e.height), y0+2); y++ {
		for x := max(0, x0-1); x < min(int(e.width), x0+2); x++ {
			if x != x0 || y != y0 {
				tile = e.field[y][x]
				if tile == Unknown {
					unknownCount++
				}
				if tile == Flag {
					flagCount++
				}
				coords = append(coords, image.Pt(x, y))
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles, coords, unknownCount, flagCount
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
