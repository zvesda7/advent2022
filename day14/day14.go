package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

type Point struct {
	x int
	y int
}

type render func(*[][]Block)

type Block byte

const (
	Empty       Block = 0
	Wall              = 1
	Sand              = 2
	Source            = 3
	OutOfBounds       = 4
)

func fillLine(board *[][]Block, a Point, b Point) {
	if a.x == b.x {
		//vertical line
		y1, y2 := a.y, b.y
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for i := y1; i <= y2; i++ {
			(*board)[a.x][i] = Wall
		}
	} else {
		//horizontal line
		x1, x2 := a.x, b.x
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for i := x1; i <= x2; i++ {
			(*board)[i][a.y] = Wall
		}
	}
}

const xStart = 400
const xEnd = 600
const yStart = 0
const yEnd = 173

const xSandSrc = 500 - xStart
const ySandSrc = 0 - yStart
const width = xEnd - xStart + 1  //92
const height = yEnd - yStart + 1 //172

func printBoard(board *[][]Block) {
	for y := 0; y < len((*board)[0]); y++ {
		for x := 0; x < len(*board); x++ {
			switch (*board)[x][y] {
			case Empty:
				fmt.Print(".")
			case Wall:
				fmt.Print("#")
			case Sand:
				fmt.Print("o")
			case Source:
				fmt.Print("+")
			}
		}
		fmt.Println()
	}
}

func getBlock(board *[][]Block, p Point) Block {
	if p.x >= 0 && p.x < len(*board) && p.y >= 0 && p.y < len((*board)[0]) {
		return (*board)[p.x][p.y]
	}
	return OutOfBounds
}

func loadBoard(paths [][]Point, addFloor bool) *[][]Block {
	board := make([][]Block, width)
	for i := 0; i < width; i++ {
		board[i] = make([]Block, height)
	}
	for _, path := range paths {
		for i := 0; i < len(path)-1; i++ {
			fillLine(&board, path[i], path[i+1])
		}
	}
	if addFloor {
		maxY := 0
		for _, path := range paths {
			for _, point := range path {
				if point.y > maxY {
					maxY = point.y
				}
			}
		}
		fillLine(&board, Point{0, maxY + 2}, Point{width - 1, maxY + 2})
	}

	return &board
}

//returns number of sand blocks rested
func runSimul(board *[][]Block, cb render) int {
	foundRest := true
	restCount := 0
	for foundRest && (*board)[xSandSrc][ySandSrc] != Sand {
		foundRest = false
		sand := Point{xSandSrc, ySandSrc}
		(*board)[sand.x][sand.y] = Sand
		for {
			(*board)[sand.x][sand.y] = Empty //remove where we were
			sand.y++
			b := getBlock(board, sand)
			if b == OutOfBounds {
				break
			} else if b != Empty {
				sand.x--
				b = getBlock(board, sand)
				if b == OutOfBounds {
					break
				} else if b != Empty {
					sand.x += 2
					b = getBlock(board, sand)
					if b == OutOfBounds {
						break
					} else if b != Empty {
						sand.x--
						sand.y--
						(*board)[sand.x][sand.y] = Sand
						foundRest = true
						restCount++
						cb(board)
						break
					}
				}
			}
			(*board)[sand.x][sand.y] = Sand
		}
	}
	return restCount
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var paths [][]Point
	for _, pathStr := range instr {
		pointsStr := strings.Split(pathStr, "->")
		var path []Point
		for _, pointStr := range pointsStr {
			xyStr := strings.Split(strings.Trim(pointStr, " "), ",")
			p := Point{}
			p.x, _ = strconv.Atoi(xyStr[0])
			p.y, _ = strconv.Atoi(xyStr[1])
			p.x -= xStart //recenter to make it easier
			p.y -= yStart
			path = append(path, p)
		}
		paths = append(paths, path)
	}

	board := loadBoard(paths, false)
	goSimulate(board)

}

func goSimulate(board *[][]Block) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 800, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	runSimul(board, func(board *[][]Block) {
		for x := 0; x < len(*board); x++ {
			for y := 0; y < len((*board)[x]); y++ {
				switch (*board)[x][y] {
				case Wall:
					rect := sdl.Rect{4 * int32(x), 4 * int32(y), 4, 4}
					surface.FillRect(&rect, 0xffff0000)
				case Sand:
					rect := sdl.Rect{4 * int32(x), 4 * int32(y), 4, 4}
					surface.FillRect(&rect, 0xff00ff00)
				}
			}
		}
		window.UpdateSurface()

		running := true
		t := time.Now().UnixMilli()
		for running && (time.Now().UnixMilli()-t) < 17 {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					println("Quit")
					running = false
					break
				}
			}
		}
	})

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

}
