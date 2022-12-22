package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const CUBE_SIZE = 24

func readLines(path string) ([]string, error) {
	var file, _ = os.ReadFile(path)
	return strings.Split(string(file), "\n"), nil
}

type Vec3 struct {
	x int
	y int
	z int
}

func addVecs(a Vec3, b Vec3) Vec3 {
	return Vec3{
		a.x + b.x,
		a.y + b.y,
		a.z + b.z,
	}
}

func hashVec(a Vec3) int {
	return a.x*10000 + a.y*100 + a.z
}

func parse(in []string) []Vec3 {
	var vecs []Vec3
	for _, s := range in {
		xyz := strings.Split(s, ",")
		//+1 hack so we can subtract one from an edge when checking adjacents
		x, _ := strconv.Atoi(xyz[0])
		y, _ := strconv.Atoi(xyz[1])
		z, _ := strconv.Atoi(xyz[2])
		vecs = append(vecs, Vec3{x + 1, y + 1, z + 1})
	}
	return vecs
}

func getAdjacent(in Vec3) []Vec3 {
	rslt := []Vec3{}
	adjVecs := []Vec3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
	}
	for _, adj := range adjVecs {
		c := addVecs(in, adj)
		if c.x >= 0 && c.x < CUBE_SIZE &&
			c.y >= 0 && c.y < CUBE_SIZE &&
			c.z >= 0 && c.z < CUBE_SIZE {
			rslt = append(rslt, c)
		}
	}
	return rslt
}

func main() {

	var instr, _ = readLines("input.txt")

	points := parse(instr)

	//build 3d voxels
	voxels := [CUBE_SIZE][CUBE_SIZE][CUBE_SIZE]byte{}

	//fill points
	for _, p := range points {
		voxels[p.x][p.y][p.z] = 1
	}

	//count unconnected sides
	unconnSides := 0
	for _, p := range points {
		for _, c := range getAdjacent(p) {
			if voxels[c.x][c.y][c.z] == 0 {
				unconnSides++
			}
		}
	}
	fmt.Println("Part 1", unconnSides)

	//BFS search to get extern
	externSides := 0
	processed := map[int]bool{}
	working := []Vec3{}
	working = append(working, Vec3{0, 0, 0})
	for len(working) > 0 {
		test := working[0]
		if !processed[hashVec(test)] {
			for _, c := range getAdjacent(test) {
				if voxels[c.x][c.y][c.z] == 0 {
					working = append(working, c)
				} else {
					externSides++
				}
			}
			processed[hashVec(test)] = true
		}
		working = working[1:]
	}

	fmt.Println("Part 2", externSides)
}
