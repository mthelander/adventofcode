package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Op func(a, b int64) int64

func add(a, b int64) int64 {
	return a + b
}

func mult(a, b int64) int64 {
	return a * b
}

func concat(a, b int64) int64 {
	numDigits := int64(math.Log10(float64(b))) + 1
	exp := int64(math.Pow(10, float64(numDigits)))
	return (a * exp) + b
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) [][]int64 {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	rg := regexp.MustCompile("[ :]+")
	var lines [][]int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := rg.Split(scanner.Text(), -1)
		line := make([]int64, len(data))

		for i, str := range data {
			i64, err := strconv.ParseInt(str, 10, 64)
			check(err)
			line[i] = i64
		}
		lines = append(lines, line)
	}
	check(scanner.Err())
	return lines
}

func solve(operators []Op, target, n int64, equation ...int64) bool {
	if len(equation) < 1 {
		return n == target
	}

	if n > target {
		return false
	}

	n2, rest := equation[0], equation[1:]

	for _, op := range operators {
		result := op(n, n2)
		if solve(operators, target, result, rest...) {
			return true
		}
	}
	return false
}

func main() {
	operators := []Op{add, mult, concat}
	filename := os.Args[1]
	sum := int64(0)
	for _, line := range readLines(filename) {
		if solve(operators, line[0], line[1], line[2:]...) {
			sum += line[0]
		}
	}
	fmt.Println(sum)
}
