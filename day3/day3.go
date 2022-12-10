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

func priority(item byte) int {
	if item >= 'a' {
		return int(item - 'a' + 1)
	} else {
		return int(item - 'A' + 1 + 26)
	}
}

func findRepeatInCompartment(backpack string) byte {
	c1 := backpack[:len(backpack)/2]
	c2 := backpack[len(backpack)/2:]
	return findMatchChar2(c1, c2)
}

func findMatchChar2(s1 string, s2 string) byte {
	for i := 0; i < len(s1); i++ {
		for k := 0; k < len(s2); k++ {
			if s1[i] == s2[k] {
				return s1[i]
			}
		}
	}
	return 0
}

func findMatchChar3(s1 string, s2 string, s3 string) byte {
	for i := 0; i < len(s1); i++ {
		for k := 0; k < len(s2); k++ {
			if s1[i] == s2[k] {
				for j := 0; j < len(s3); j++ {
					if s1[i] == s3[j] {
						return s1[i]
					}
				}
			}
		}
	}
	return 0
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//part 1
	sum1 := 0
	for _, x := range instr {
		sum1 += priority(findRepeatInCompartment(x))
	}
	fmt.Println("Part1", sum1)

	//part 2
	sum2 := 0
	for i := 0; i < len(instr); i += 3 {
		sum2 += priority(findMatchChar3(instr[i], instr[i+1], instr[i+2]))
	}
	fmt.Println("Part2", sum2)
}
