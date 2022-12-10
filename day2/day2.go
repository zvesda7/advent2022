package main

import (
	"fmt"
	"os"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(strings.TrimSpace(string(file)), "\n"), nil
}

var scoreMap = map[string]int{
	"A X": 1 + 3,
	"A Y": 2 + 6,
	"A Z": 3 + 0,
	"B X": 1 + 0,
	"B Y": 2 + 3,
	"B Z": 3 + 6,
	"C X": 1 + 6,
	"C Y": 2 + 0,
	"C Z": 3 + 3,
}

var scoreMap2 = map[string]int{
	"A X": 3 + 0,
	"A Y": 1 + 3,
	"A Z": 2 + 6,
	"B X": 1 + 0,
	"B Y": 2 + 3,
	"B Z": 3 + 6,
	"C X": 2 + 0,
	"C Y": 3 + 3,
	"C Z": 1 + 6,
}

func main() {
	fmt.Println("hello")
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	sum1, sum2 := 0, 0
	for _, x := range instr {
		sum1 += scoreMap[x]
		sum2 += scoreMap2[x]
	}
	fmt.Println(sum1, sum2)
}
