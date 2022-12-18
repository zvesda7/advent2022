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
	height       int
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
	pf.cells = append(pf.cells, [9]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1})
	for y := 0; y < height-1; y++ {
		pf.cells = append(pf.cells, [9]uint8{1, 0, 0, 0, 0, 0, 0, 0, 1})
	}
	pf.topFilledRow = 0
	pf.height = height
	return &pf
}

func (pf *PlayField) calcY(y int) int {
	return y & (pf.height - 1)
}

func (pf *PlayField) getTopFilledRow() int {
	for {
		clear := true
		for x := 1; x < len(pf.cells[0])-1; x++ {
			if pf.cells[pf.calcY(pf.topFilledRow+1)][x] != 0 {
				clear = false
			}
		}
		if clear {
			break
		}
		pf.topFilledRow++
		zeroY := pf.calcY(pf.topFilledRow + 7)
		for x := 1; x < len(pf.cells[0])-1; x++ {
			pf.cells[zeroY][x] = 0
		}
	}

	return pf.topFilledRow
}

func (pf *PlayField) checkCanMove(p *Piece, x int, y int) bool {
	for dx := 0; dx < 4; dx++ {
		for dy := 0; dy < 4; dy++ {
			if p.bitmap[dy][dx] == 1 && pf.cells[pf.calcY(y-dy)][x+dx] == 1 {
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
				pf.cells[pf.calcY(y-dy)][x+dx] = 1
			}
		}
	}
}

func printBottom10(pf *PlayField, p Piece, px int, py int) {
	for y := 10; y >= 0; y-- {
		for x := 0; x < len(pf.cells[0]); x++ {
			if pf.cells[y][x] == 1 {
				fmt.Print("#")
			} else {
				if x >= px && x < (px+4) && y > (py-4) && y <= py {
					if p.bitmap[py-y][x-px] == 1 {
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

func printTop10(pf *PlayField) {
	for y := pf.topFilledRow + 1; y >= pf.topFilledRow-21; y-- {
		for x := 0; x < len(pf.cells[0]); x++ {
			if pf.cells[pf.calcY(y)][x] == 1 {
				fmt.Print("#")
			} else {

				fmt.Print(".")

			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func calcLarge(pattern string, numPieces int) int {
	leftStartPos := 2
	bottomStartPos := 3

	windNext := makeWindGen(pattern)
	pieces := buildPieces()
	pieceNext := makeNextPiece(pieces)
	pf := buildPlayfield(512)
	for i := 0; i < numPieces; i++ {
		p := pieceNext()
		px := leftStartPos + 1
		py := pf.getTopFilledRow() + bottomStartPos + p.height
		rested := false
		for !rested {
			wind := windNext()
			if pf.checkCanMove(p, px+wind, py) {
				px += wind
			}
			//printBottom10(pf, *p, px, py)
			if pf.checkCanMove(p, px, py-1) {
				py--
			} else {
				rested = true
				pf.setPiece(p, px, py)
			}
			//printBottom10(pf, *p, px, py)
		}
	}
	fmt.Println("Part1", pf.getTopFilledRow())
	printTop10(pf)
	return pf.getTopFilledRow()
}

func main() {

	var instr, _ = readLines("test.txt")

	//pieces dropped repeat every 5pieces*10091winds = 50455 times
	//the pattern at top of stack after 50455 matches pattern at 50455*7 perfectly
	//formula is ((x-50455)%50455)+50455, take this number and solve that stack.
	repeatN := 5 * len(instr[0])
	fmt.Println(repeatN)

	x1 := calcLarge(instr[0], repeatN)
	x8 := calcLarge(instr[0], repeatN*8)

	x15 := calcLarge(instr[0], repeatN*15)
	x22 := calcLarge(instr[0], repeatN*22)
	x7 := x15 - x8
	//x15_s := x1 + 2*x7
	fmt.Println(x1, x8, x7, x8-x1, x15-x8, x22-x15)

	n := 1000000000000

	x := x1
	n -= repeatN

	x += n / (repeatN * 7) * x7
	n = n % (repeatN * 7)

	x += calcLarge(instr[0], repeatN+n) - x1
	n = 0
	fmt.Println("Part 2", x)

}
