package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

type obj struct {
	hasVal   bool
	val      int
	children []*obj
}

type objSlice []*obj

func (a objSlice) Len() int           { return len(a) }
func (a objSlice) Less(i, j int) bool { return compareObj(a[i], a[j]) == -1 }
func (a objSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func initObj() *obj {
	return &obj{
		children: make([]*obj, 0),
	}
}

//not needed in solution but helpful when troubleshooting
func (o obj) String() string {
	if len(o.children) > 0 {
		a := "["
		for i := 0; i < len(o.children)-1; i++ {
			a = a + o.children[i].String() + ","
		}
		a = a + o.children[len(o.children)-1].String() + "]"
		return a
	} else if o.hasVal {
		return strconv.Itoa(o.val)
	} else {
		return ""
	}
}

func parseObj(s string) *obj {
	root := &obj{}

	var stack []*obj
	var curObj *obj = root
	for _, c := range s {
		switch c {
		case '[':
			stack = append(stack, curObj)
			curObj = initObj()
		case ']':
			parent := stack[len(stack)-1]
			parent.children = append(parent.children, curObj)
			curObj = parent
			stack = stack[:len(stack)-1]
		case ',':
			parent := stack[len(stack)-1]
			parent.children = append(parent.children, curObj)
			curObj = initObj()
		default:
			curObj.val = curObj.val*10 + int(c-'0')
			curObj.hasVal = true
		}
	}
	return root
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func compareInt(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func compareObj(a *obj, b *obj) int {
	if a.hasVal && b.hasVal {
		return compareInt(a.val, b.val)
	} else if a.hasVal && !b.hasVal {
		n := initObj()
		n.children = append(n.children, a)
		return compareObj(n, b)
	} else if b.hasVal && !a.hasVal {
		n := initObj()
		n.children = append(n.children, b)
		return compareObj(a, n)
	} else {
		count := min(len(a.children), len(b.children))
		for i := 0; i < count; i++ {
			if cmp := compareObj(a.children[i], b.children[i]); cmp != 0 {
				return cmp
			}
		}
		return compareInt(len(a.children), len(b.children))
	}
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//parse packets
	var allPackets []*obj
	for i := 0; i < len(instr); i += 3 {
		allPackets = append(allPackets, parseObj(instr[i]))
		allPackets = append(allPackets, parseObj(instr[i+1]))
	}

	sumIndexes := 0
	for i := 0; i < len(allPackets); i += 2 {
		if compareObj(allPackets[i], allPackets[i+1]) == -1 {
			sumIndexes += i/2 + 1
		}
	}
	fmt.Println("Part 1", sumIndexes)

	div1 := parseObj("[[2]]")
	div2 := parseObj("[[6]]")
	allPackets = append(allPackets, div1)
	allPackets = append(allPackets, div2)
	sort.Sort(objSlice(allPackets))

	i1, i2 := 0, 0
	for i, p := range allPackets {
		if p == div1 {
			i1 = i + 1
		} else if p == div2 {
			i2 = i + 1
		}
	}
	fmt.Println("Part 2", i1*i2)
}
