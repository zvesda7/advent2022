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

//column height, string of last top 10 rows
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
			if pf.checkCanMove(p, px, py-1) {
				py--
			} else {
				rested = true
				pf.setPiece(p, px, py)
			}
		}
	}
	return pf.getTopFilledRow()
}

func main() {

	var instr, _ = readLines("input.txt")

	//pieces dropped repeat every 5pieces*10091winds = 50455 times
	//the pattern at top of stack after 50455*1 matches pattern at 50455*349, and 50455*697
	//found by search 50455*1 through 50455*1000

	repeatN := 5 * len(instr[0])
	//dupliTracker := map[string][]int{}
	//strs := calcLargeMultiple(instr[0], repeatN*1000, repeatN)
	//for i, s := range strs {
	//	dupliTracker[s] = append(dupliTracker[s], i)
	//}
	//for _, lst := range dupliTracker {
	//	for _, k := range lst {
	//		fmt.Print(k, ",")
	//	}
	//	fmt.Println()
	//}

	repeatM := 348

	x1 := calcLarge(instr[0], repeatN)
	xM := calcLarge(instr[0], repeatN*(repeatM+1))

	xM2 := calcLarge(instr[0], repeatN*(repeatM*2+1))
	xMD := xM2 - xM

	//fmt.Println(x1, x8, x7, x8-x1, x15-x8, x22-x15)

	n := 1000000000000

	x := x1
	n -= repeatN

	x += n / (repeatN * repeatM) * xMD
	n = n % (repeatN * repeatM)

	x += calcLarge(instr[0], repeatN+n) - x1
	n = 0
	fmt.Println("Part 2", x)

}
