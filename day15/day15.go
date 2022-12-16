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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func remove(slice []Range, s int) []Range {
	return append(slice[:s], slice[s+1:]...)
}

type Sensor struct {
	x       int
	y       int
	beaconX int
	beaconY int
}

type Range struct {
	x     int
	count int
}

func parseSensors(instr []string) []Sensor {
	var sensors []Sensor
	getNums := regexp.MustCompile("[-0-9]+")
	for _, line := range instr {
		numStrings := getNums.FindAllString(line, -1)
		s := Sensor{}
		s.x, _ = strconv.Atoi(numStrings[0])
		s.y, _ = strconv.Atoi(numStrings[1])
		s.beaconX, _ = strconv.Atoi(numStrings[2])
		s.beaconY, _ = strconv.Atoi(numStrings[3])
		sensors = append(sensors, s)
	}
	return sensors
}

func calcManhat(s Sensor) int {
	return abs(s.x-s.beaconX) + abs(s.y-s.beaconY)
}

func reduceRanges(r []Range) []Range {
	//combine ranges until they no longer overlap, then sum them
	sort.Slice(r, func(i, j int) bool {
		return r[i].x < r[j].x
	})
	for i := 0; i < len(r); i++ {

		for j := i + 1; j < len(r); {
			if r[j].x <= r[i].x+r[i].count {
				newCount := r[j].x - r[i].x + r[j].count
				if newCount > r[i].count {
					r[i].count = newCount
				}
				r = remove(r, j)
			} else {
				break
			}
		}
	}
	return r
}

func findRangesAtRow(sensors []Sensor, y int) []Range {
	ranges := make([]Range, 0)
	for _, s := range sensors {
		d := calcManhat(s)
		d2y := abs(s.y - y)
		if d2y <= d {
			dx := d - d2y
			ranges = append(ranges, Range{s.x - dx, dx*2 + 1})
		}
	}
	return reduceRanges(ranges)
}

func boxRanges(r []Range, boxLen int) []Range {
	for i := len(r) - 1; i >= 0; i-- {
		if r[i].x < 0 {
			r[i].count += r[i].x
			r[i].x = 0
		}
		if r[i].x+r[i].count > boxLen {
			r[i].count += boxLen - (r[i].x + r[i].count)
		}
		if r[i].count < 0 {
			r = remove(r, i)
		}
	}
	return r
}

func main() {
	const TestRow = 2000000
	const BoxLen = 4000001
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	sensors := parseSensors(instr)

	beaconXs := make(map[int]bool)
	for _, s := range sensors {
		if s.beaconY == TestRow {
			beaconXs[s.beaconX] = true
		}
	}
	ranges := findRangesAtRow(sensors, TestRow)
	allPoints := 0
	for _, x := range ranges {
		allPoints += x.count
	}
	fmt.Println("Part 1", allPoints-len(beaconXs))

	for y := 0; y < BoxLen; y++ {
		ranges := findRangesAtRow(sensors, y)
		ranges = boxRanges(ranges, BoxLen)
		if len(ranges) == 2 {
			x := ranges[1].x - 1
			fmt.Println("Part 2", 4000000*x+y)
			break
		}
	}
}
