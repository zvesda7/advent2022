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

type Node struct {
	name      string
	flowRate  int
	links     []*Node
	linkNames []string
}

type TreeNode struct {
	parent   *TreeNode
	children []*TreeNode
	node     *Node
	score    int
	steps    int
	todo     []*Node
	keep     bool
}

func parseNodes(instr []string) map[string]*Node {
	nodes := make(map[string]*Node)
	getNums := regexp.MustCompile("[-0-9]+")
	getNodeNames := regexp.MustCompile("[A-Z][A-Z]")
	for _, line := range instr {
		numStrings := getNums.FindAllString(line, -1)
		nameStrings := getNodeNames.FindAllString(line, -1)
		n := Node{}
		n.name = nameStrings[0]
		n.flowRate, _ = strconv.Atoi(numStrings[0])
		for i := 1; i < len(nameStrings); i++ {
			n.linkNames = append(n.linkNames, nameStrings[i])
		}
		nodes[n.name] = &n
	}

	for _, n := range nodes {
		for _, name := range n.linkNames {
			n.links = append(n.links, nodes[name])
		}
	}

	return nodes
}

func getValvesToOpen(nodes map[string]*Node) []*Node {
	var toOpen []*Node
	for _, n := range nodes {
		if n.flowRate > 0 {
			toOpen = append(toOpen, n)
		}
	}
	return toOpen
}
func findShortestPath(startN *Node, endN *Node) int {
	//BFS
	dists := map[string]int{}
	visited := map[string]bool{}

	var workingSet []*Node
	workingSet = append(workingSet, startN)
	dists[startN.name] = 0
	for len(workingSet) > 0 {
		var currNode *Node = workingSet[0]
		workingSet = workingSet[1:]
		for _, n := range currNode.links {
			if !visited[n.name] {
				dists[n.name] = dists[currNode.name] + 1
				workingSet = append(workingSet, n)
				visited[n.name] = true
				if n == endN {
					return dists[endN.name]
				}
			}
		}
	}
	return dists[endN.name]
}

func calcTotalFlow(step int, node *Node) int {
	if step > 30 {
		return 0
	}
	return (30 - step) * node.flowRate
}

func remove(nodes []*Node, index int) []*Node {
	newArr := make([]*Node, len(nodes)-1)
	for i, v := range nodes {
		if i < index {
			newArr[i] = v
		} else if i > index {
			newArr[i-1] = v
		}
	}
	return newArr
}

func calcTree(t *TreeNode, depth int, maxDepth int, lastLayer *[]*TreeNode) {
	if depth > maxDepth {
		*lastLayer = append(*lastLayer, t)
		return
	}
	if t.todo != nil {
		for i, x := range t.todo {
			d := findShortestPath(t.node, x) + 1 //plus 1 for turning the valve on
			newt := TreeNode{}
			newt.parent = t
			newt.node = x
			newt.steps = t.steps + d
			newt.score = t.score + calcTotalFlow(newt.steps, x)
			newt.todo = remove(t.todo, i)
			t.children = append(t.children, &newt)
		}
		t.todo = nil
	}
	for _, x := range t.children {
		calcTree(x, depth+1, maxDepth, lastLayer)
	}

	return
}

func markKeep(t *TreeNode) {
	t.keep = true
	if t.parent != nil {
		markKeep(t.parent)
	}
}

func pruneRecurv(t *TreeNode) {
	t.keep = false //reset for next pass
	for i := len(t.children) - 1; i >= 0; i-- {
		if !t.children[i].keep {
			t.children[i] = t.children[len(t.children)-1]
			t.children = t.children[:len(t.children)-1]
		} else {
			pruneRecurv(t.children[i])
		}
	}
}

func pruneTree(root *TreeNode, lastLayer *[]*TreeNode, curDepth int, maxDepth int) {
	reduceFactor := maxDepth - curDepth + 1
	(*lastLayer) = (*lastLayer)[len(*lastLayer)-len(*lastLayer)/reduceFactor:]
	//fmt.Println("highest", (*lastLayer)[len(*lastLayer)-1].steps)
	//fmt.Println(reduceFactor, len(*lastLayer))
	for _, x := range *lastLayer {
		markKeep(x)
	}
	pruneRecurv(root)
}

func main() {
	initialDepth := 4
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	nodes := parseNodes(instr)
	valves := getValvesToOpen(nodes)

	tree := TreeNode{}
	tree.todo = valves
	tree.node = nodes["AA"]
	maxDepth := len(valves)
	for depth := initialDepth; depth <= maxDepth; depth++ {
		lastLayer := []*TreeNode{}
		calcTree(&tree, 1, depth, &lastLayer)
		sort.Slice(lastLayer, func(i, j int) bool {
			return (lastLayer)[i].score < (lastLayer)[j].score
		})
		fmt.Println(lastLayer[len(lastLayer)-1].score)
		if depth != maxDepth {
			pruneTree(&tree, &lastLayer, depth, maxDepth)
		}
	}

}
