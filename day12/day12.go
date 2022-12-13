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

type Node struct {
	x    int
	y    int
	z    int
	dist int
}

func (n Node) GetHash() int {
	return n.x*10000 + n.y
}

func getAllAdj(grid *[][]*Node, n *Node) []*Node {
	var nodes []*Node
	var deltas = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, delta := range deltas {
		x, y := n.x+delta[0], n.y+delta[1]
		if y >= 0 && y < len(*grid) && x >= 0 && x < len((*grid)[y]) {
			nt := (*grid)[y][x]
			if nt.z >= (n.z - 1) { //reverse logic, we only traverse end to start to make part 2 easier
				nodes = append(nodes, nt)
			}
		}
	}
	return nodes
}

func findShortestPath(grid *[][]*Node, startN *Node, endN *Node, endOnAnyOfSameZAsEnd bool) int {
	//BFS
	var workingSet []*Node
	workingSet = append(workingSet, startN)
	visited := map[int]bool{}
	for len(workingSet) > 0 {
		var currNode *Node = workingSet[0]
		workingSet = workingSet[1:]
		adjNodes := getAllAdj(grid, currNode)
		for _, n := range adjNodes {
			if !visited[n.GetHash()] {
				n.dist = currNode.dist + 1
				workingSet = append(workingSet, n)
				visited[n.GetHash()] = true
			}
		}
		if endOnAnyOfSameZAsEnd && currNode.z == endN.z {
			return currNode.dist
		}
	}
	return endN.dist
}

func loadGrid(instr []string) ([][]*Node, *Node, *Node) {
	var startN *Node
	var endN *Node
	grid := make([][]*Node, len(instr))
	for y, row := range instr {
		grid[y] = make([]*Node, len(row))
		for x, cell := range row {
			z := int(cell)
			thisNode := Node{x, y, z, 0}
			grid[thisNode.y][thisNode.x] = &thisNode
			if cell == 'S' {
				startN = &thisNode
				startN.dist = 0
				startN.z = int('a')
			} else if cell == 'E' {
				endN = &thisNode
				endN.z = int('z')
			}
		}
	}
	return grid, startN, endN
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	grid, startN, endN := loadGrid(instr)
	part1 := findShortestPath(&grid, endN, startN, false)
	fmt.Println("part1", part1)

	//walk backward, looking for any 'a'
	grid, startN, endN = loadGrid(instr)
	part2 := findShortestPath(&grid, endN, startN, true)
	fmt.Println("part2", part2)
}
