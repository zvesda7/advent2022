package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

type Sensor struct {
	x       int
	y       int
	beaconX int
	beaconY int
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

func main() {
	var instr, err = readLines("test.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	sensors := parseSensors(instr)
	fmt.Println(sensors)
}
