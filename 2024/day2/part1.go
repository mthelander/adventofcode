package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
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

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func readLines(path string) [][]int {
	file, err := os.Open(path)
	check(err)
	defer file.Close()
	rg := regexp.MustCompile("\\s+")

	var matrix [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := rg.Split(scanner.Text(), -1)
		matrix = append(matrix, Map(data, parseInt))
	}
	check(scanner.Err())

	return matrix
}

func difference(a, b int) int {
	return int(math.Abs(float64(a - b)))
}

func isSafe(list []int) bool {
	return isMonotonic(list) && isGradual(list)
}

func isMonotonic(list []int) bool {
	if len(list) <= 1 {
		return true
	}

	increasing := list[1] > list[2]

	for i := range len(list) - 1 {
		if (increasing && list[i] < list[i+1]) || (!increasing && list[i] > list[i+1]) {
			return false
		}
	}

	return true
}

func isGradual(list []int) bool {
	if len(list) <= 1 {
		return true
	}
	for i := range len(list) - 1 {
		delta := difference(list[i], list[i+1])
		if delta < 1 || delta > 3 {
			return false
		}
	}

	return true
}

func main() {
	fname := os.Args[1]
	matrix := readLines(fname)

	numSafe := 0

	for _, row := range matrix {
		if isSafe(row) {
			numSafe += 1
		}
	}

	fmt.Println("Num Safe: ", numSafe)
}
