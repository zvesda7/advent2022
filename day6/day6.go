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

func noRepeats(s string) bool {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				return false
			}
		}
	}
	return true
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	data := instr[0]

	for i := 4; i <= len(data); i++ {
		last4 := data[i-4 : i]
		if noRepeats(last4) {
			fmt.Println("part1", i)
			break
		}
	}
	for i := 14; i <= len(data); i++ {
		last14 := data[i-14 : i]
		if noRepeats(last14) {
			fmt.Println("part2", i)
			break
		}
	}

}
