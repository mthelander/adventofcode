package main

import (
	"bufio"
	"fmt"
	"math"
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
	return left, right
}

func main() {
	fname := os.Args[1]
	sum := 0
	left, right := readLines(fname)
	slices.Sort(left)
	slices.Sort(right)

	for i := range left {
		delta := float64(left[i] - right[i])
		sum += int(math.Abs(delta))
	}

	fmt.Println("Sum: ", sum)
}
