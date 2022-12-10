package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	fmt.Println("hello")
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var elves []int
	var curElfCals int
	for _, x := range instr {
		if x == "" {
			elves = append(elves, curElfCals)
			curElfCals = 0
		} else {
			var pint, err = strconv.Atoi(x)
			if err != nil {
				os.Exit(0)
			}
			curElfCals = curElfCals + pint
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(elves)))
	var total int
	for i := 0; i < 3; i++ {
		fmt.Println("elf", i, elves[i])
		total += elves[i]
	}
	fmt.Println("total", total)
}
