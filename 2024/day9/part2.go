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

func willitfit(diskmap *[]string, pos, size int) bool {
	if pos >= len(*diskmap) {
		return false
	}
	if size > 0 {
		n := (*diskmap)[pos]
		return n == "." && willitfit(diskmap, pos+1, size-1)
	}
	return true
}

func defrag(diskmap *[]string, start, end int) bool {
	if start <= end {
		if (*diskmap)[start] != "." {
			return defrag(diskmap, start+1, end)
		}
		if (*diskmap)[end] == "." {
			return defrag(diskmap, start, end-1)
		}
		// start == free space, end == last id of chunk
		size := sizeof(diskmap, end)
		// end - size == start of chunk to move
		for i := start; i < end; i++ {
			if willitfit(diskmap, i, size) {
				swapchunk(diskmap, i, end, size)
				return defrag(diskmap, start, end-size)
			}
		}
		// No place for this file block
		return defrag(diskmap, start, end-size)
	}
	return true
}

func swapchunk(diskmap *[]string, i, j, size int) {
	if size > 0 {
		swap(diskmap, i, j)
		swapchunk(diskmap, i+1, j-1, size-1)
	}
}

func sizeof(diskmap *[]string, n int) int {
	id := (*diskmap)[n]
	size := 0
	for i := n; (*diskmap)[i] == id && id != "."; i-- {
		size++
	}
	return size
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
