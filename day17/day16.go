package main

import (
	"fmt"
	"os"
	"strings"
)

type Piece struct {
	bitmap [4][4]uint8 //y,x
	height int
	width  int
}

type PlayField struct {
	cells        [][9]uint8
	topFilledRow int
	yoffset      int
}

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

func makeWindGen(s string) func() int {
	i := -1
	return func() int {
		if i++; i >= len(s) {
			i = 0
		}
		if s[i] == '>' {
			return 1
		} else {
			return -1
		}
	}
}

func makeNextPiece(ps []*Piece) func() *Piece {
	i := -1
	return func() *Piece {
		if i++; i >= len(ps) {
			i = 0
		}
		return ps[i]
	}
}

func buildPieces() []*Piece {
	var ps []*Piece
	ps = append(ps, &Piece{
		[4][4]uint8{
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		1,
		4,
	})

	ps = append(ps, &Piece{
		[4][4]uint8{
			{0, 1, 0, 0},
			{1, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
		3,
		3,
	})
	ps = append(ps, &Piece{
		[4][4]uint8{
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
		},
		3,
		3,
	})
	ps = append(ps, &Piece{
		[4][4]uint8{
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
		},
		4,
		1,
	})
	ps = append(ps, &Piece{
		[4][4]uint8{
			{1, 1, 0, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		2,
		2,
	})
	return ps
}

func buildPlayfield(height int) *PlayField {
	pf := PlayField{}
	for y := 0; y < height-1; y++ {
		pf.cells = append(pf.cells, [9]uint8{1, 0, 0, 0, 0, 0, 0, 0, 1})
	}
	pf.cells = append(pf.cells, [9]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1})
	pf.topFilledRow = height - 1
	return &pf
}

func (pf *PlayField) getTopFilledRow() int {
	topClearRow := pf.topFilledRow - 1
	for {
		clear := true
		for x := 1; x < len(pf.cells[0])-1; x++ {
			if pf.cells[topClearRow][x] != 0 {
				clear = false
			}
		}
		if clear {
			break
		}
		topClearRow--
	}
	pf.topFilledRow = topClearRow + 1
	return pf.topFilledRow
}

func (pf *PlayField) checkCanMove(p *Piece, x int, y int) bool {
	for dx := 0; dx < 4; dx++ {
		for dy := 0; dy < 4; dy++ {
			if p.bitmap[dy][dx] == 1 && pf.cells[y+dy][x+dx] == 1 {
				return false
			}
		}
	}
	return true
}

func (pf *PlayField) setPiece(p *Piece, x int, y int) {
	for dx := 0; dx < 4; dx++ {
		for dy := 0; dy < 4; dy++ {
			if p.bitmap[dy][dx] == 1 {
				pf.cells[y+dy][x+dx] = 1
			}
		}
	}
}

func printBottom10(pf *PlayField, p Piece, px int, py int) {
	for y := len(pf.cells) - 18; y < len(pf.cells); y++ {
		for x := 0; x < len(pf.cells[0]); x++ {
			if pf.cells[y][x] == 1 {
				fmt.Print("#")
			} else {
				if x >= px && x < (px+4) && y >= py && y < (py+4) {
					if p.bitmap[y-py][x-px] == 1 {
						fmt.Print("@")
					} else {
						fmt.Print(".")
					}
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()

}

func main() {
	leftStartPos := 2
	bottomStartPos := 3
	numPieces := 2022
	var instr, _ = readLines("input.txt")

	windNext := makeWindGen(instr[0])
	pieces := buildPieces()
	pieceNext := makeNextPiece(pieces)
	pf := buildPlayfield(numPieces * 4)
	for i := 0; i < numPieces; i++ {
		p := pieceNext()
		px := leftStartPos + 1
		py := pf.getTopFilledRow() - bottomStartPos - p.height
		rested := false
		for !rested {
			//printBottom10(pf, p, px, py)
			wind := windNext()
			if pf.checkCanMove(p, px+wind, py) {
				px += wind
			}
			if pf.checkCanMove(p, px, py+1) {
				py++
			} else {
				rested = true
				pf.setPiece(p, px, py)
			}
		}
	}
	fmt.Println("Part1", len(pf.cells)-pf.getTopFilledRow()-1)

}
