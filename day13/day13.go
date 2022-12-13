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

type obj struct {
	val  int
	vals []*obj
}

func (o *obj) isArray() bool {
	return len(o.vals) > 0
}
func (o obj) String() string {
	if o.isArray() {
		a := "["
		for i := 0; i < len(o.vals)-1; i++ {
			a = a + o.vals[i].String() + ","
		}
		a = a + o.vals[len(o.vals)-1].String() + "]"
		return a
	} else {
		return strconv.Itoa(o.val)
	}
}

type packetPair struct {
	a *obj
	b *obj
}

func parseObj(s string) *obj {
	root := &obj{}

	var stack []*obj
	var curObj *obj = root
	for _, c := range s {
		switch c {
		case '[':
			stack = append(stack, curObj)
			curObj = &obj{}
		case ']':
			parent := stack[len(stack)-1]
			parent.vals = append(parent.vals, curObj)
			curObj = parent
			stack = stack[:len(stack)-1]
		case ',':
			parent := stack[len(stack)-1]
			parent.vals = append(parent.vals, curObj)
			curObj = &obj{}
		default:
			curObj.val = curObj.val*10 + int(c-'0')
		}
		fmt.Println(len(root.vals))
	}
	return root
}

func main() {
	//var instr, err = readLines("test.txt")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(0)
	//}

	o := parseObj("[[4,4],4,4,4]")
	fmt.Println(o)

	//parse packets
	//var packetPairs []packetPair
	//for i := 0; i < len(instr); i += 3 {
	//	p := packetPair{}
	//	p.a = parseObj(instr[i])
	//	p.b = parseObj(instr[i+1])
	//	packetPairs = append(packetPairs, p)
	//}

}
