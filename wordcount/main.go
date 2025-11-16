//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func processFile(filename string) map[string]int {
	data, err := os.ReadFile(filename)
	errCheck(err)

	linesCount := make(map[string]int)

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		_, exists := linesCount[line]
		if !exists {
			linesCount[line] = 1
		} else {
			linesCount[line] += 1
		}
	}

	return linesCount
}

func updateMap(map1 map[string]int, map2 map[string]int) {
	for m2k, m2v := range map2 {
		_, exists := map1[m2k]
		if !exists {
			map1[m2k] = m2v
		} else {
			map1[m2k] += m2v
		}
	}
}

func printLinesCount(linesCountMap map[string]int) {
	for k, v := range linesCountMap {
		if v <= 1 {
			continue
		}
		fmt.Printf("%d\t%s\n", v, k)
	}
}

func main() {
	filenames := os.Args[1:]

	resMap := make(map[string]int)
	for _, filename := range filenames {
		linesCount := processFile(filename)
		updateMap(resMap, linesCount)
	}
	printLinesCount(resMap)
}
