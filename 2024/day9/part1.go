package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runeToi(r rune) int {
	i, err := strconv.Atoi(string(r))
	check(err)
	return i
}

func readLines(path string) []int {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := make([]int, 0)
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		for _, char := range scanner.Text() {
			data = append(data, runeToi(char))
		}
	}
	check(scanner.Err())
	return data
}

func expandfile(result *[]string, id, head int, rest ...int) {
	for range head {
		*result = append(*result, strconv.Itoa(id))
	}
	if len(rest) > 0 {
		expandfree(result, id, rest[0], rest[1:]...)
	}
}

func expandfree(result *[]string, id, head int, rest ...int) {
	for range head {
		*result = append(*result, ".")
	}
	if len(rest) > 0 {
		expandfile(result, id+1, rest[0], rest[1:]...)
	}
}

func expand(diskmap []int) *[]string {
	result := make([]string, 0)
	expandfile(&result, 0, diskmap[0], diskmap[1:]...)
	return &result
}

func swap(diskmap *[]string, a, b int) {
	(*diskmap)[a], (*diskmap)[b] = (*diskmap)[b], (*diskmap)[a]
}

func defrag(diskmap *[]string, start, end int) bool {
	if start <= end {
		if (*diskmap)[start] != "." {
			return defrag(diskmap, start+1, end)
		}
		if (*diskmap)[end] == "." {
			return defrag(diskmap, start, end-1)
		}
		swap(diskmap, start, end)
		return defrag(diskmap, start+1, end-1)
	}
	return true
}

func checksum(diskmap *[]string, pos int) int64 {
	if pos >= len(*diskmap) {
		return 0
	}
	s := (*diskmap)[pos]
	if s == "." {
		return checksum(diskmap, pos+1)
	}
	i, err := strconv.ParseInt(s, 10, 64)
	check(err)
	return (i * int64(pos)) + checksum(diskmap, pos+1)
}

func main() {
	diskmap := readLines(os.Args[1])
	expanded := expand(diskmap)
	defrag(expanded, 0, len(*expanded)-1)
	fmt.Printf("%v\n", checksum(expanded, 0))
}
