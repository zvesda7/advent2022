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

const inf = 9999

type Node struct {
	x       int
	y       int
	z       int
	dist    int
	visited bool
}

func pHash(x, y int) int {
	return x*10000 + y
}

func getZ(instr *[]string, n *Node) byte {
	return byte((*instr)[n.y][n.x])
}

func getAdj(allNodes *map[int]*Node, n *Node, dir int, reverseHeightRule bool) *Node {
	x := n.x
	y := n.y
	switch dir {
	case 0:
		x++
	case 1:
		x--
	case 2:
		y++
	case 3:
		y--
	}
	if nt, found := (*allNodes)[pHash(x, y)]; found {
		if reverseHeightRule {
			if nt.z >= (n.z - 1) {
				return nt
			}
		} else {
			if nt.z <= (n.z + 1) {
				return nt
			}
		}
	}
	return nil
}

func getAllAdj(allNodes *map[int]*Node, n *Node, reverseHeightRule bool) []*Node {
	var nodes []*Node
	for dir := 0; dir < 4; dir++ {
		if nt := getAdj(allNodes, n, dir, reverseHeightRule); nt != nil {
			nodes = append(nodes, nt)
		}
	}
	return nodes
}

func findShortestPath(allNodes *map[int]*Node, startN *Node, endN *Node, endOnAnyOfSameZAsEnd bool, reverseHeightRule bool) int {
	for (*allNodes)[pHash(endN.x, endN.y)].dist == inf {
		var currNode *Node
		for _, n := range *allNodes {
			if n.visited == false && (currNode == nil || n.dist < currNode.dist) {
				currNode = n
			}
		}

		adjNodes := getAllAdj(allNodes, currNode, reverseHeightRule)
		nextDist := currNode.dist + 1
		for _, n := range adjNodes {
			if nextDist < n.dist {
				n.dist = nextDist
			}
		}
		currNode.visited = true
		if endOnAnyOfSameZAsEnd && currNode.z == endN.z {
			return currNode.dist
		}
		delete(*allNodes, pHash(currNode.x, currNode.y))
		//fmt.Println(currNode)
	}
	return endN.dist
}

func loadNodes(instr []string) (map[int]*Node, *Node, *Node) {
	var startN *Node
	var endN *Node
	allNodes := make(map[int]*Node) //key is node hash, value is distance to reach it.
	for y, row := range instr {
		for x, cell := range row {
			z := int(cell)
			thisNode := Node{x, y, z, inf, false}
			allNodes[pHash(thisNode.x, thisNode.y)] = &thisNode
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
	return allNodes, startN, endN
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	allNodes, startN, endN := loadNodes(instr)
	part1 := findShortestPath(&allNodes, startN, endN, false, false)
	fmt.Println("part1", part1)

	//walk backward, looking for any 'a'
	allNodes2, startN2, endN2 := loadNodes(instr)
	startN2.dist = inf
	endN2.dist = 0
	part2 := findShortestPath(&allNodes2, endN2, startN2, true, true)
	fmt.Println("part2", part2)
}
