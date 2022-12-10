package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

var t = 0
var x = 1
var sumSignal = 0

func nextCycle() {
	t++
	if (t+20)%40 == 0 {
		sumSignal += t * x
	}
	xBeam := (t - 1) % 40
	if xBeam == 0 {
		fmt.Print("\n")
	}
	if xBeam >= (x-1) && xBeam <= (x+1) {
		fmt.Print("#")
	} else {
		fmt.Print(".")
	}
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for _, ins := range instr {
		nextCycle()
		spl := strings.Split(ins, " ")
		if spl[0] == "addx" {
			nextCycle()
			dx, _ := strconv.Atoi(spl[1])
			x += dx
		}
	}

	fmt.Println("Part 1", sumSignal)

}
