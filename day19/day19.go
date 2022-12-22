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

type Blueprint struct {
	nth      int
	ore_ore  int
	clay_ore int
	obs_ore  int
	obs_clay int
	geo_ore  int
	geo_obs  int
}

func parseBlueprints(instr []string) []Blueprint {
	var bps []Blueprint
	getNums := regexp.MustCompile("[-0-9]+")
	for _, line := range instr {
		numStrings := getNums.FindAllString(line, -1)
		s := Blueprint{}
		s.nth, _ = strconv.Atoi(numStrings[0])
		s.ore_ore, _ = strconv.Atoi(numStrings[1])
		s.clay_ore, _ = strconv.Atoi(numStrings[2])
		s.obs_ore, _ = strconv.Atoi(numStrings[3])
		s.obs_clay, _ = strconv.Atoi(numStrings[4])
		s.geo_ore, _ = strconv.Atoi(numStrings[5])
		s.geo_obs, _ = strconv.Atoi(numStrings[6])
		bps = append(bps, s)
	}
	return bps
}

func main() {

	var instr, _ = readLines("test.txt")

	bps := parseBlueprints(instr)
	fmt.Println(bps)
}
