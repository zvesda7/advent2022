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

type file struct {
	name string
	size int
}

type dir struct {
	name   string
	files  []file
	dirs   []*dir
	parent *dir
}

func printDir(d *dir, level int) {
	pad := strings.Repeat(" ", level*2)
	println(pad + "- " + d.name + " (dir)")
	for _, subDir := range d.dirs {
		printDir(subDir, level+1)
	}
	for _, file := range d.files {
		println(pad + "  - " + file.name + " (file, size " + fmt.Sprint(file.size) + ")")
	}
}

func sumFiles(d *dir) int {
	total := 0
	for _, sub := range d.dirs {
		total += sumFiles(sub)
	}
	for _, file := range d.files {
		total += file.size
	}
	return total
}

func listDirs(d *dir) []*dir {
	var dirs []*dir
	dirs = append(dirs, d)
	for _, sub := range d.dirs {
		dirs = append(dirs, listDirs(sub)...)
	}
	return dirs
}

func main() {
	var instr, err = readLines("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	root := dir{name: "/"}
	currDir := &root

	for _, line := range instr {
		if line == "$ cd /" {
			currDir = &root
		} else if line == "$ cd .." {
			currDir = currDir.parent
		} else if len(line) >= 4 && line[:4] == "$ cd" {
			newDirName := line[5:]
			for _, subdir := range currDir.dirs {
				if subdir.name == newDirName {
					currDir = subdir
				}
			}
		} else if line == "$ ls" {
			//do nothing
		} else if len(line) >= 4 && line[:4] == "dir " {
			newDir := dir{
				name:   line[4:],
				parent: currDir,
			}
			currDir.dirs = append(currDir.dirs, &newDir)
		} else {
			//file
			parts := strings.Split(line, " ")
			newFile := file{}
			newFile.name = parts[1]
			newFile.size, _ = strconv.Atoi(parts[0])
			currDir.files = append(currDir.files, newFile)
		}
	}
	printDir(&root, 0)

	usedSpace := sumFiles(&root)
	freeSpace := 70000000 - usedSpace
	needSpace := 30000000 - freeSpace

	allDirs := listDirs(&root)
	total := 0
	smallestGreaterThanNeeded := 999999999
	for _, d := range allDirs {
		size := sumFiles(d)
		if size <= 100000 {
			total += size
		}
		if size >= needSpace {
			fmt.Println(size)
		}
		if size >= needSpace && size < smallestGreaterThanNeeded {
			smallestGreaterThanNeeded = size
		}
	}
	fmt.Println("Part1", total)
	fmt.Println("Part2", smallestGreaterThanNeeded)
}
