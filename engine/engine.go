package engine

import (
	"bytes"
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
	GameLoop()
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
	e.LeftClick(int(e.width/2), int(e.height/2)) // Start somewhere
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

// GameLoop handles game logic and communication
func (e *engine) GameLoop() {
	var i, j int
	for {
		didSomething := false
		e.UpdateField()
		e.PrintField()
		for i = 0; i < int(e.width); i++ {
			for j = 0; j < int(e.height); j++ {
				tile := e.field[i][j]
				log.Println(i, j, tile)
				if tile < 1 || tile > 8 {
					continue
				}
				tiles, coords, unknownCount, flagCount := e.getNeighbours(i, j)
				log.Println(tilesString(tiles))
				// Marking flags
				if unknownCount == int(tile)-flagCount {
					for k := 0; k < len(coords); k++ {
						c := coords[k]
						t := e.field[c.X][c.Y]
						if t == Unknown {
							e.field[c.X][c.Y] = Flag
							flagCount++
							log.Println("Setting flag at", c.X, c.Y)
							e.RightClick(c.Y, c.X)
							didSomething = true
						}
					}
				}
				// Clicking on safe unknowns
				if unknownCount > 0 && int(tile) == flagCount {
					for k := 0; k < len(coords); k++ {
						c := coords[k]
						t := e.field[c.X][c.Y]
						if t == Unknown {
							log.Println("Clicking on", c.X, c.Y)
							e.LeftClick(c.Y, c.X)
							didSomething = true
						}
					}
				}
			}
		}
		if !didSomething {
			log.Println("🌀 Cannot decide what to do..")
			break
		}
	}
}

func (e engine) getNeighbours(x, y int) ([]Tile, []image.Point, int, int) {
	var tiles = []Tile{}
	var coords = []image.Point{}
	var unknownCount, flagCount int
	var tile Tile
	for i := max(0, x-1); i < min(int(e.width), x+2); i++ {
		for j := max(0, y-1); j < min(int(e.height), y+2); j++ {
			if i != x || j != y {
				tile = e.field[i][j]
				if tile == Unknown {
					unknownCount++
				}
				if tile == Flag {
					flagCount++
				}
				coords = append(coords, image.Pt(i, j))
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
