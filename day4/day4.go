package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(strings.TrimSpace(string(file)), "\n"), nil
}

type elfPair struct {
	startA int
	endA   int
	startB int
	endB   int
}

func newElfPair(textLine string) *elfPair {
	textPairs := strings.Split(textLine, ",")
	textElfA := strings.Split(textPairs[0], "-")
	textElfB := strings.Split(textPairs[1], "-")

	p := elfPair{}
	p.startA, _ = strconv.Atoi(textElfA[0])
	p.endA, _ = strconv.Atoi(textElfA[1])
	p.startB, _ = strconv.Atoi(textElfB[0])
	p.endB, _ = strconv.Atoi(textElfB[1])
	return &p
}

func fullyContained(e *elfPair) bool {
	if e.startA >= e.startB && e.endA <= e.endB {
		return true
	}
	if e.startB >= e.startA && e.endB <= e.endA {
		return true
	}
	return false
}

func anyOverlap(e *elfPair) bool {
	if e.startA >= e.startB && e.startA <= e.endB {
		return true
	}
	if e.startB >= e.startA && e.startB <= e.endA {
		return true
	}
	return false
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//part 1
	sum1, sum2 := 0, 0
	for _, x := range instr {
		p := newElfPair(x)
		if fullyContained(p) {
			sum1 += 1
		}
		if anyOverlap(p) {
			sum2 += 1
		}
	}
	fmt.Println("Part1", sum1, "Part2", sum2)

}
