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
	nth       int
	name      string
	flowRate  int
	links     []*Node
	linkNames []string
}

type TreeNode struct {
	parent   *TreeNode
	children []*TreeNode
	node1    *Node
	node2    *Node
	score    int
	steps1   int
	steps2   int
	todo     []*Node
	keep     bool
}

func parseNodes(instr []string) map[string]*Node {
	nodes := make(map[string]*Node)
	getNums := regexp.MustCompile("[-0-9]+")
	getNodeNames := regexp.MustCompile("[A-Z][A-Z]")
	for i, line := range instr {
		numStrings := getNums.FindAllString(line, -1)
		nameStrings := getNodeNames.FindAllString(line, -1)
		n := Node{}
		n.nth = i + 1
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

var pathPreCalcs map[int]int

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
	if step > 26 {
		return 0
	}
	return (26 - step) * node.flowRate
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

type NodePair struct {
	node1 *Node
	node2 *Node
}

func chooseTwo(list []*Node) []NodePair {
	var nps []NodePair
	for i := 0; i < len(list)-1; i++ {
		for j := i + 1; j < len(list); j++ {
			np := NodePair{
				list[i],
				list[j],
			}
			nps = append(nps, np)
			np2 := NodePair{
				list[j],
				list[i],
			}
			nps = append(nps, np2)
		}
	}
	return nps
}

func remove2(nodes []*Node, n1 *Node, n2 *Node) []*Node {
	newArr := make([]*Node, len(nodes)-2)
	k := 0
	for _, v := range nodes {
		if v.nth != n1.nth && v.nth != n2.nth {
			newArr[k] = v
			k++
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
		pairs := chooseTwo(t.todo)
		for _, pair := range pairs {

			newt := TreeNode{}
			newt.parent = t
			d := pathPreCalcs[t.node1.nth*1000+pair.node1.nth] + 1 //plus 1 for turning the valve on
			newt.node1 = pair.node1
			newt.steps1 = t.steps1 + d

			d2 := pathPreCalcs[t.node2.nth*1000+pair.node2.nth] + 1 //plus 1 for turning the valve on
			newt.node2 = pair.node2
			newt.steps2 = t.steps2 + d2

			newt.score = t.score + calcTotalFlow(newt.steps1, pair.node1) + calcTotalFlow(newt.steps2, pair.node2)
			newt.todo = remove2(t.todo, pair.node1, pair.node2)
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
	reduceFactor := 5 * curDepth
	(*lastLayer) = (*lastLayer)[len(*lastLayer)-len(*lastLayer)/reduceFactor:]
	//fmt.Println("highest", (*lastLayer)[len(*lastLayer)-1].steps)
	//fmt.Println(reduceFactor, len(*lastLayer))
	for _, x := range *lastLayer {
		markKeep(x)
	}
	pruneRecurv(root)
}

func main() {
	initialDepth := 2
	var instr, err = readLines("test.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	nodes := parseNodes(instr)
	valves := getValvesToOpen(nodes)

	//sort.Slice(valves, func(i, j int) bool {
	//	return (valves)[i].name < (valves)[j].name
	//})

	//precalc distances
	pathPreCalcs = make(map[int]int)
	valvesWithAA := append(valves, nodes["AA"])
	for _, x := range valvesWithAA {
		for _, y := range valvesWithAA {
			pathPreCalcs[x.nth*1000+y.nth] = findShortestPath(x, y)
		}
	}

	tree := TreeNode{}
	tree.todo = valves
	tree.node1 = nodes["AA"]
	tree.node2 = nodes["AA"]
	maxDepth := len(valves) / 2
	for depth := initialDepth; depth <= maxDepth; depth++ {
		lastLayer := []*TreeNode{}
		calcTree(&tree, 1, depth, &lastLayer)
		sort.Slice(lastLayer, func(i, j int) bool {
			return (lastLayer)[i].score < (lastLayer)[j].score
		})
		fmt.Println(len(lastLayer), lastLayer[len(lastLayer)-1].score)
		if depth != maxDepth {
			pruneTree(&tree, &lastLayer, depth, maxDepth)
		} else {
			printStack(lastLayer[len(lastLayer)-1])
			//for i := 0; i < len(lastLayer); i++ {
			//fmt.Println(lastLayer[i].score)
			//printStack(lastLayer[i])
			//}
		}
	}

}

func printStack(lowestNode *TreeNode) {
	fmt.Println(lowestNode.node1.name, lowestNode.node2.name)
	if lowestNode.parent != nil {
		printStack(lowestNode.parent)
	}
}

//EE CC
//HH BB
//DD JJ
//AA AA
