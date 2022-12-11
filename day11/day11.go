package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

type operationFunc func(int64) int64

type monkey struct {
	numInspects    int
	items          []int64
	operation      operationFunc
	testDivisor    int64
	ifTrueMonkeyI  int
	ifFalseMonkeyI int
}

func loadMonkeys(instr []string) []monkey {
	var monkeys []monkey
	getNums := regexp.MustCompile("[0-9]+")
	operatorIndex := 23
	operandIndex := 25

	for index := 0; index < len(instr); index++ {
		m := monkey{}

		index++ // starting items
		startingItems := getNums.FindAllString(instr[index], -1)
		m.items = make([]int64, len(startingItems))
		for i, s := range startingItems {
			m.items[i], _ = strconv.ParseInt(s, 10, 64)
		}

		index++ //operation
		operandStr := strings.TrimSpace(instr[index])[operandIndex-2:]
		operand, numericErr := strconv.ParseInt(operandStr, 10, 64)
		switch instr[index][operatorIndex] {
		case '+':
			if numericErr == nil {
				fmt.Println("+", operand)
				m.operation = func(x int64) int64 { return x + operand }
			} else if operandStr == "old" {
				fmt.Println("+", "old")
				m.operation = func(x int64) int64 { return x + x }
			}
		case '*':
			if numericErr == nil {
				fmt.Println("*", operand)
				m.operation = func(x int64) int64 { return x * operand }
			} else if operandStr == "old" {
				fmt.Println("*", "old")
				m.operation = func(x int64) int64 { return x * x }
			}
		}
		index++ //test
		testDivisor, _ := strconv.ParseInt(getNums.FindAllString(instr[index], -1)[0], 10, 64)
		m.testDivisor = testDivisor
		index++ //test true
		monkeyTrue, _ := strconv.Atoi(getNums.FindAllString(instr[index], -1)[0])
		m.ifTrueMonkeyI = monkeyTrue
		index++ //test false
		monkeyFalse, _ := strconv.Atoi(getNums.FindAllString(instr[index], -1)[0])
		m.ifFalseMonkeyI = monkeyFalse
		index++ //space

		monkeys = append(monkeys, m)
	}

	return monkeys
}

func calcMonkeyBusiness(monkeys []monkey, numRounds int, div3Worry bool) int64 {
	for round := 0; round < numRounds; round++ {
		for m := 0; m < len(monkeys); m++ {
			monkey := &monkeys[m]
			for len(monkey.items) > 0 {
				monkey.numInspects++
				item := monkey.items[0]
				monkey.items = monkey.items[1:]
				item = monkey.operation(item)
				if div3Worry {
					item = item / 3
				}
				if (item % monkey.testDivisor) == 0 {
					monkeys[monkey.ifTrueMonkeyI].items = append(monkeys[monkey.ifTrueMonkeyI].items, item)
				} else {
					monkeys[monkey.ifFalseMonkeyI].items = append(monkeys[monkey.ifFalseMonkeyI].items, item)
				}
			}
		}
	}

	numInspects := make([]int, len(monkeys))
	for i, m := range monkeys {
		numInspects[i] = m.numInspects
	}
	fmt.Println(numInspects)
	sort.Ints((numInspects))

	fmt.Println("m1", numInspects[len(numInspects)-1])
	fmt.Println("m2", numInspects[len(numInspects)-2])

	return int64(numInspects[len(numInspects)-1]) * int64(numInspects[len(numInspects)-2])
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	monkeys := loadMonkeys(instr)
	part1 := calcMonkeyBusiness(monkeys, 20, true)
	fmt.Println("Part 1", part1)

	monkeys2 := loadMonkeys(instr)
	part2 := calcMonkeyBusiness(monkeys2, 20, false)
	fmt.Println("Part 2", part2)
}
