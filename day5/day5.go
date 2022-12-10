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

type stack []byte

func (s *stack) push(b byte) {
	*s = append(*s, b)
}
func (s *stack) pop() byte {
	n := len(*s) - 1
	x := (*s)[n]
	*s = (*s)[:n]
	return x
}
func (s *stack) top() byte {
	n := len(*s) - 1
	x := (*s)[n]
	return x
}

type instruct struct {
	count int
	src   int
	dest  int
}

func combineTops(stacks []stack) string {
	cmb := ""
	for _, s := range stacks {
		cmb += string(rune(s.top()))
	}
	return cmb
}

func parsefile(instr []string) ([]stack, []instruct) {
	numStacks := (len(instr[0]) + 1) / 4
	stacks := make([]stack, numStacks)

	var bottomIndex int
	for i, s := range instr {
		if s[0:2] == " 1" {
			bottomIndex = i
			break
		}
	}
	for i := bottomIndex - 1; i >= 0; i-- {
		for j := 0; j < numStacks; j++ {
			pos := j*4 + 1
			if instr[i][pos] != ' ' {
				stacks[j].push(instr[i][pos])
			}
		}
	}

	var instructs []instruct
	for i := bottomIndex + 2; i < len(instr); i++ {
		splits := strings.Split(instr[i], " ")
		ins := instruct{}
		ins.count, _ = strconv.Atoi(splits[1])
		ins.src, _ = strconv.Atoi(splits[3])
		ins.dest, _ = strconv.Atoi(splits[5])

		instructs = append(instructs, ins)
	}

	return stacks, instructs
}

func runInstructs(stacks []stack, instructs []instruct) []stack {
	for _, x := range instructs {
		for i := 0; i < x.count; i++ {
			stacks[x.dest-1].push(stacks[x.src-1].pop())
		}
	}
	return stacks
}

func runInstructs2(stacks []stack, instructs []instruct) []stack {
	for _, x := range instructs {
		tempStack := stack{}
		for i := 0; i < x.count; i++ {
			tempStack.push(stacks[x.src-1].pop())
		}
		for i := 0; i < x.count; i++ {
			stacks[x.dest-1].push(tempStack.pop())
		}
	}
	return stacks
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	stacks, instructs := parsefile(instr)
	stacks = runInstructs(stacks, instructs)
	fmt.Println("Part1", combineTops(stacks))

	stacks, instructs = parsefile(instr)
	stacks = runInstructs2(stacks, instructs)
	fmt.Println("Part2", combineTops(stacks))
}
