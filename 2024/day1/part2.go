package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseInt(data string) int {
	i, err := strconv.Atoi(data)
	check(err)
	return i
}

func readLines(path string) ([]int, []int) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()
	rg := regexp.MustCompile("\\s+")

	var left, right []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := rg.Split(scanner.Text(), 2)
		left = append(left, parseInt(data[0]))
		right = append(right, parseInt(data[1]))
	}
	check(scanner.Err())

	slices.Sort(left)
	slices.Sort(right)

	return left, right
}

func frequencyMap(list []int) map[int]int {
	var lookup map[int]int
	lookup = make(map[int]int)

	for _, e := range list {
		lookup[e] += 1
	}

	return lookup
}

func main() {
	fname := os.Args[1]
	sum := 0
	left, right := readLines(fname)
	freq := frequencyMap(right)

	for _, v := range left {
		sum += v * freq[v]
	}

	fmt.Println("Sum: ", sum)
}
