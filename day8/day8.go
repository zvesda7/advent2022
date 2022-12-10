package main

import (
	"fmt"
	"os"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	yCount := len(instr)
	xCount := len(instr[0])

	vismap := make([][]byte, yCount)
	for i := range vismap {
		vismap[i] = make([]byte, xCount)
	}

	//horizontal,
	for y := 1; y < yCount-1; y++ {
		//left to right
		h := instr[y][0]
		for x := 1; x < xCount-1; x++ {
			th := instr[y][x]
			if th > h {
				vismap[y][x] = 1
				h = th
			}
		}
		//right to left
		h = instr[y][xCount-1]
		for x := xCount - 2; x >= 1; x-- {
			th := instr[y][x]
			if th > h {
				vismap[y][x] = 1
				h = th
			}
		}
	}

	//vertical
	for x := 1; x < xCount-1; x++ {
		//top to bottom
		h := instr[0][x]
		for y := 1; y < yCount-1; y++ {
			th := instr[y][x]
			if th > h {
				vismap[y][x] = 1
				h = th
			}
		}
		//bottom to top
		h = instr[yCount-1][x]
		for y := yCount - 2; y >= 1; y-- {
			th := instr[y][x]
			if th > h {
				vismap[y][x] = 1
				h = th
			}
		}
	}

	//sum
	totalVis := yCount*2 + (xCount-2)*2 //four sides

	for y := 0; y < yCount; y++ {
		for x := 0; x < xCount; x++ {
			if vismap[y][x] > 0 {
				totalVis++
			}
		}
	}
	fmt.Println("Part 1", totalVis)

	maxScore := 0
	for y := 1; y < yCount-1; y++ {
		for x := 1; x < xCount-1; x++ {
			h := instr[y][x]

			score1 := 0 //up
			for y2 := y - 1; y2 >= 0; y2-- {
				score1++
				if instr[y2][x] >= h {
					break
				}
			}
			score2 := 0 //down
			for y2 := y + 1; y2 < yCount; y2++ {
				score2++
				if instr[y2][x] >= h {
					break
				}
			}
			score3 := 0 //left
			for x2 := x - 1; x2 >= 0; x2-- {
				score3++
				if instr[y][x2] >= h {
					break
				}
			}
			score4 := 0 //right
			for x2 := x + 1; x2 < xCount; x2++ {
				score4++
				if instr[y][x2] >= h {
					break
				}
			}

			score := score1 * score2 * score3 * score4
			if score > maxScore {
				maxScore = score
			}
		}
	}

	fmt.Println("Part 2", maxScore)
}
