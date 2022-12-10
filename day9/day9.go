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

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

type seg struct {
	x     int
	y     int
	child *seg
}

func (s *seg) move(d string) {
	if d == "U" {
		s.y--
	} else if d == "D" {
		s.y++
	} else if d == "L" {
		s.x--
	} else if d == "R" {
		s.x++
	}
	if s.child != nil {
		s.child.follow(s)
	}
}

func (s *seg) follow(parent *seg) {
	movedx, movedy := false, false
	if parent.x-s.x > 1 {
		s.x++
		movedx = true
	} else if s.x-parent.x > 1 {
		s.x--
		movedx = true
	}
	if parent.y-s.y > 1 {
		s.y++
		movedy = true
	} else if s.y-parent.y > 1 {
		s.y--
		movedy = true
	}

	if movedy && !movedx && parent.x-s.x > 0 {
		s.x++
	} else if movedy && !movedx && s.x-parent.x > 0 {
		s.x--
	}
	if movedx && !movedy && parent.y-s.y > 0 {
		s.y++
	} else if movedx && !movedy && s.y-parent.y > 0 {
		s.y--
	}

	if s.child != nil {
		s.child.follow(s)
	}
}

func (s *seg) count() int {
	if s.child == nil {
		return 1
	} else {
		return 1 + s.child.count()
	}
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	head1 := seg{}
	tail1 := seg{}
	head1.child = &tail1
	uniquexy1 := make(map[int]bool)
	uniquexy1[tail1.x*10000+tail1.y] = true

	head2 := seg{}
	var tail2 *seg
	tail2 = &head2
	for i := 0; i < 9; i++ {
		t := seg{}
		tail2.child = &t
		tail2 = &t
	}
	uniquexy2 := make(map[int]bool)
	uniquexy2[tail2.x*1000000+tail2.y] = true

	for _, ins := range instr {
		both := strings.Split(ins, " ")
		dir := both[0]
		count, _ := strconv.Atoi(both[1])
		for i := 0; i < count; i++ {
			head1.move(dir)
			uniquexy1[tail1.x*10000+tail1.y] = true

			head2.move(dir)
			uniquexy2[tail2.x*1000000+tail2.y] = true
		}
	}
	fmt.Println("Part 1", len(uniquexy1))
	fmt.Println("Part 2", len(uniquexy2))
	fmt.Println(head2.count())
}
